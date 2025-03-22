# Go-Site

My custom site and static site generator.

## TODO

- Move to different Markdown parser or find a way to have
  text/template strings persist after transformation

- Deploy to a hosting service
  - e.g. Netlify, Github pages, Cloudflare, ...
   - https://github.com/emad-elsaid/xlog/blob/master/tutorials/Create%20your%20own%20digital%20garden%20on%20Github.md
- Add build badges to README
- Move build tooling to it's own repo

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
