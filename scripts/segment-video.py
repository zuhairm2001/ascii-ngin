#!/usr/bin/env python3
import argparse
import io
import subprocess
import sys
import tempfile
from pathlib import Path

from PIL import Image
from rembg import new_session, remove
from tqdm import tqdm
import imageio_ffmpeg


def run(cmd):
    subprocess.run(cmd, check=True)


def ffmpeg_exe():
    return imageio_ffmpeg.get_ffmpeg_exe()


def get_input_fps(ffmpeg, input_path):
    cmd = [ffmpeg, "-i", str(input_path)]
    proc = subprocess.run(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE, text=True
    )
    stderr = proc.stderr
    for line in stderr.splitlines():
        if "Video:" in line and "fps" in line:
            parts = line.split(",")
            for part in parts:
                part = part.strip()
                if part.endswith(" fps"):
                    try:
                        return float(part.replace(" fps", ""))
                    except ValueError:
                        continue
    return None


def extract_frames(ffmpeg, input_path, frames_dir, fps):
    frames_dir.mkdir(parents=True, exist_ok=True)
    cmd = [ffmpeg, "-i", str(input_path)]
    if fps:
        cmd += ["-r", str(fps)]
    cmd += [str(frames_dir / "%05d.png")]
    run(cmd)


def remove_background(frames_dir, cut_dir, model_name):
    cut_dir.mkdir(parents=True, exist_ok=True)
    frames = sorted(frames_dir.glob("*.png"))
    session = new_session(model_name)
    for frame in tqdm(frames, desc="Removing background", unit="frame"):
        with Image.open(frame) as im:
            out = remove(im, session=session)
            if isinstance(out, Image.Image):
                out_im = out
            else:
                out_im = Image.open(io.BytesIO(out))
            out_im.save(cut_dir / frame.name)


def build_mov(ffmpeg, cut_dir, output_path, fps):
    cmd = [
        ffmpeg,
        "-framerate",
        str(fps),
        "-i",
        str(cut_dir / "%05d.png"),
        "-c:v",
        "prores_ks",
        "-profile:v",
        "4",
        "-pix_fmt",
        "yuva444p10le",
        str(output_path),
    ]
    run(cmd)


def build_webm(ffmpeg, cut_dir, output_path, fps):
    cmd = [
        ffmpeg,
        "-framerate",
        str(fps),
        "-i",
        str(cut_dir / "%05d.png"),
        "-c:v",
        "libvpx-vp9",
        "-pix_fmt",
        "yuva420p",
        "-b:v",
        "0",
        "-crf",
        "30",
        str(output_path),
    ]
    run(cmd)


def parse_args():
    parser = argparse.ArgumentParser(
        description="Segment subject and output transparent video."
    )
    parser.add_argument(
        "-i",
        "--input",
        required=True,
        help="Path to input video file",
    )
    parser.add_argument(
        "--format",
        choices=["mov", "webm", "both"],
        default="both",
        help="Output format (default: both)",
    )
    parser.add_argument(
        "--fps",
        type=float,
        default=None,
        help="Force output FPS (default: source FPS)",
    )
    parser.add_argument(
        "--model",
        default="u2net_human_seg",
        help="rembg model name (default: u2net_human_seg)",
    )
    parser.add_argument(
        "--keep-frames",
        action="store_true",
        help="Keep extracted frames",
    )
    return parser.parse_args()


def main():
    args = parse_args()
    input_path = Path(args.input)
    if not input_path.exists():
        print(f"Input file not found: {input_path}", file=sys.stderr)
        sys.exit(1)

    ffmpeg = ffmpeg_exe()
    fps = args.fps or get_input_fps(ffmpeg, input_path)
    if not fps:
        fps = 30.0
        print("Could not detect FPS, defaulting to 30.", file=sys.stderr)

    if args.keep_frames:
        frames_dir = Path("frames")
        cut_dir = Path("cut")
        extract_frames(ffmpeg, input_path, frames_dir, args.fps)
        remove_background(frames_dir, cut_dir, args.model)
    else:
        with tempfile.TemporaryDirectory() as tmp:
            tmp_path = Path(tmp)
            frames_dir = tmp_path / "frames"
            cut_dir = tmp_path / "cut"
            extract_frames(ffmpeg, input_path, frames_dir, args.fps)
            remove_background(frames_dir, cut_dir, args.model)
            if args.format in ("mov", "both"):
                output_mov = input_path.with_name(f"{input_path.stem}-segmented.mov")
                build_mov(ffmpeg, cut_dir, output_mov, fps)
            if args.format in ("webm", "both"):
                output_webm = input_path.with_name(f"{input_path.stem}-segmented.webm")
                build_webm(ffmpeg, cut_dir, output_webm, fps)
            return

    if args.format in ("mov", "both"):
        output_mov = input_path.with_name(f"{input_path.stem}-segmented.mov")
        build_mov(ffmpeg, cut_dir, output_mov, fps)
    if args.format in ("webm", "both"):
        output_webm = input_path.with_name(f"{input_path.stem}-segmented.webm")
        build_webm(ffmpeg, cut_dir, output_webm, fps)


if __name__ == "__main__":
    main()
