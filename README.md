# Go-Site

My custom site and static site generator.

[![Build and Deploy](https://github.com/NDoolan360/Go-Site/actions/workflows/deploy.yml/badge.svg?branch=main)](https://github.com/NDoolan360/Go-Site/actions/workflows/deploy.yml)

## Environment Variables

- ENV: Environment (dev, prod)
- GITHUB_USERNAME: Github username
- GITHUB_TOKEN: Github token
- CULTS3D_USERNAME: Cults3D username
- CULTS3D_API_KEY: Cults3D API key
- BGG_GEEKLIST: BoardGameGeek geeklist id

## Docker

```bash
docker-compose up
```

## Roadmap

- [x] Add Docker support
- [x] Add Github Actions deployment
- [x] Add Github Actions automated testing
- [ ] Add RSS feed
- [ ] Add article/project pagination
- [ ] Add article/project search
- [ ] Add cron to refresh data from external sources
- [ ] Build a custom markdown parser
