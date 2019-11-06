package copy

import (
	"context"
	"io"
)

// Helper function for copy from dst to src
// Similar with io.copyBuffer except it actually break loop on io.EOF (in this telnet Client it is significant)
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {

	for {
		nr, er := src.Read(buf)

		// break on any read error (and also on EOF)
		if er != nil {
			err = er
			break
		}

		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

	}
	return written, err
}

// Helper to Run copyBuffer-ing process in go-routine
// Params
//  - context for canceling
//  - Reader from where will read
//  - Writer where will write
//  - channel of error that could be happened
func RunCopier(ctx context.Context, reader io.Reader, writer io.Writer, result chan<- error) {
	go func() {

		// buffer for copyBuffer
		buf := make([]byte, 4096)

	LOOP:
		for {
			select {
			case <-ctx.Done():
				result <- nil
				break LOOP
			default:
				// Writer <- Reader
				_, err := copyBuffer(writer, reader, buf)
				if err != nil {
					result <- err
					break LOOP
				}
			}
		}
	}()
}
