package client // import "github.com/docker/docker/client"

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types/image"
)

// ImageLoad loads an image in the docker host from the client host.
// It's up to the caller to close the io.ReadCloser in the
// ImageLoadResponse returned by this function.
//
// Platform is an optional parameter that specifies the platform to load from
// the provided multi-platform image. This is only has effect if the input image
// is a multi-platform image.
func (cli *Client) ImageLoad(ctx context.Context, input io.Reader, opts image.LoadOptions) (image.LoadResponse, error) {
	v := url.Values{}
	v.Set("quiet", "0")
	if opts.Quiet {
		v.Set("quiet", "1")
	}
	if opts.Platform != nil {
		if err := cli.NewVersionError(ctx, "1.48", "platform"); err != nil {
			return image.LoadResponse{}, err
		}

		p, err := json.Marshal(*opts.Platform)
		if err != nil {
			return image.LoadResponse{}, err
		}
		v.Set("platform", string(p))
	}

	resp, err := cli.postRaw(ctx, "/images/load", v, input, http.Header{
		"Content-Type": {"application/x-tar"},
	})
	if err != nil {
		return image.LoadResponse{}, err
	}
	return image.LoadResponse{
		Body: resp.body,
		JSON: resp.header.Get("Content-Type") == "application/json",
	}, nil
}
