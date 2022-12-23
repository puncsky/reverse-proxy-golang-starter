# reverse-proxy-golang-starter

A handy but powerful light-weight reverse proxy starter template with Golang.

## How to use it?

Update `main.go` to forward traffic as you wish. Now `/blog` points to `https://example.com` and `/` points to `https://tianpan.co`:

```golang
	cfgs := []ForwardCfg{
		{"/blog", "https://example.com", false},
		{"/", "https://tianpan.co", true},
	}
```

and then run the project

```bash
make install
make dev
```

or build to binary

```bash
make build
PORT=4321 ./app
```

## Licence

MIT
