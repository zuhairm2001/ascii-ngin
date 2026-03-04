package video

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

const ndjsonMaxLineSize = 10 * 1024 * 1024

func WriteNDJSON(frames <-chan FrameRecord, w io.Writer, meta MetaRecord) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if meta.Type == "" {
		meta.Type = "meta"
	}
	if err := encoder.Encode(meta); err != nil {
		return fmt.Errorf("failed to write meta record: %w", err)
	}

	for frame := range frames {
		if err := encoder.Encode(frame); err != nil {
			return fmt.Errorf("failed to write frame record: %w", err)
		}
	}

	return nil
}

func ReadNDJSON(r io.Reader) (<-chan FrameRecord, <-chan error) {
	frames := make(chan FrameRecord)
	errs := make(chan error, 1)

	go func() {
		defer close(frames)
		defer close(errs)

		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), ndjsonMaxLineSize)

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var meta MetaRecord
			if err := json.Unmarshal(line, &meta); err == nil && meta.Type == "meta" {
				continue
			}

			var frame FrameRecord
			if err := json.Unmarshal(line, &frame); err != nil {
				errs <- fmt.Errorf("failed to parse frame record: %w", err)
				return
			}
			frames <- frame
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("failed to read ndjson: %w", err)
		}
	}()

	return frames, errs
}
