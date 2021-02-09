package clair

import (
	"context"

	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	name "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/quay/claircore"
)

func rt(ref string) (http.RoundTripper, error) {

	r, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	repo := r.Context()
	key := repo.String()
	rtMu.Lock()
	defer rtMu.Unlock()

	if v, ok := rtMap[key]; ok {
		return v, nil
	}

	auth, err := authn.DefaultKeychain.Resolve(repo)
	if err != nil {
		return nil, err
	}

	rt, err := transport.New(repo.Registry, auth, http.DefaultTransport, []string{repo.Scope("pull")})
	if err != nil {
		return nil, err
	}

	rtMap[key] = rt
	return rt, nil
}

func inspect(ctx context.Context, r string) (*claircore.Manifest, error) {

	rt, err := rt(r)
	if err != nil {
		return nil, err
	}

	ref, err := name.ParseReference(r)
	if err != nil {
		return nil, err
	}
	desc, err := remote.Get(ref, remote.WithTransport(rt))
	if err != nil {
		return nil, err
	}
	img, err := desc.Image()
	if err != nil {
		return nil, err
	}
	dig, err := img.Digest()
	if err != nil {
		return nil, err
	}
	ccd, err := claircore.ParseDigest(dig.String())
	if err != nil {
		return nil, err
	}
	out := claircore.Manifest{Hash: ccd}

	ls, err := img.Layers()
	if err != nil {
		return nil, err
	}

	repo := ref.Context()
	rURL := url.URL{
		Scheme: repo.Scheme(),
		Host:   repo.RegistryStr(),
	}
	c := http.Client{
		Transport: rt,
	}

	for _, l := range ls {
		d, err := l.Digest()
		if err != nil {
			return nil, err
		}
		ccd, err := claircore.ParseDigest(d.String())
		if err != nil {
			return nil, err
		}
		u, err := rURL.Parse(path.Join("/", "v2", strings.TrimPrefix(repo.RepositoryStr(), repo.RegistryStr()), "blobs", d.String()))
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
		if err != nil {
			return nil, err
		}
		res, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		err = res.Body.Close()
		if err != nil {
			return nil, err
		}

		res.Request.Header.Del("User-Agent")
		out.Layers = append(out.Layers, &claircore.Layer{
			Hash:    ccd,
			URI:     res.Request.URL.String(),
			Headers: res.Request.Header,
		})
	}

	return &out, nil
}
