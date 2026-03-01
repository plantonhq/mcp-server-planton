package domains

import (
	"fmt"
	"io"
	"strings"

	"google.golang.org/grpc"
)

// DrainStream reads messages from a gRPC server-streaming RPC until EOF,
// context cancellation, or maxEntries is reached. Each message is formatted
// by the caller-supplied format function, keeping proto-specific imports out
// of this package.
//
// Partial results are returned on context deadline (expected for running
// jobs whose streams do not close). The caller is responsible for creating
// a bounded context before opening the stream.
func DrainStream[T any](
	stream grpc.ServerStreamingClient[T],
	maxEntries int,
	format func(*T) string,
) (string, int, error) {
	var buf strings.Builder
	count := 0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if count > 0 {
				break
			}
			return "", 0, err
		}

		buf.WriteString(format(msg))
		buf.WriteByte('\n')
		count++

		if maxEntries > 0 && count >= maxEntries {
			fmt.Fprintf(&buf, "\n--- output truncated at %d entries ---\n", maxEntries)
			break
		}
	}

	if count == 0 {
		return "No log entries available.", 0, nil
	}

	var out strings.Builder
	fmt.Fprintf(&out, "Collected %d log entries.\n\n", count)
	out.WriteString(buf.String())
	return out.String(), count, nil
}
