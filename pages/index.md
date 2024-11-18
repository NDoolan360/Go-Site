---
Title: Nathan Doolan
Description: Software Engineer. 3D printing and boardgame nerd.
  Shipping code at [Kaluza](https://kaluza.com), tinkering on
  [side projects](#projects), and learning a little along the way.
Template: /index_template.html
StyleSheets:
  - /static/styles/utils/index.css
  - /static/styles/utils/posts.css
  - /static/styles/utils/projects.css
PostLimit: 3
ProjectLimit: 10
---

<section id="intro">

![Profile Image](static/images/profile.webp "Picture of me at the Strawberry Field in Liverpool")

# {{.Page.Title}}

{{.Page.Description}}

</section>
<section id="links">

- [![Github Logo](static/images/logos/github.svg) Github](https://github.com/NDoolan360)
- [![LinkedIn Logo](static/images/logos/linkedin.svg) LinkedIn](https://www.linkedin.com/in/nathan-doolan-835a13171)
- [![Discord Logo](static/images/logos/discord.svg) Discord](https://discord.com/users/nothindoin)
- [![Cults3d Logo](static/images/logos/cults3d.svg) Cults3D](https://cults3d.com/en/users/ND360)
- [![Bgg Logo](static/images/logos/bgg.svg) BGG](https://github.com/NDoolan360)
- [![Email Icon](static/images/email.svg) Email](mailto:mail@doolan.dev)

</section>

## [Posts](/blog.html) {#posts}

{{ range $i, $_ := .Site.BlogPosts }}
  {{ if lt $i $.Page.PostLimit }}
- <a class="post-card" href="{{ .URL }}">

  ### {{ .Title }} <time datetime="{{ .Created }}">{{ .Created }}</time>
  {{ .Description }}

  </a>
  {{ end }}
{{ end }}

## [Projects](/projects.html) {#projects}

<ul>
{{ range $i, $_ := .Site.Projects }}
  {{ if lt $i $.Page.ProjectLimit }}
<li class="project-card">
<a class="project-card-link" href="{{ .URL }}"></a>
<div class="project-card-content">

### {{ .Title }} ![Logo]({{ .Logo }})
{{if .Description}}{{ .Description }}{{end}}
{{ range .Tags }}
- {{ . }}
{{ end }}

</div>
<img class="project-card-image" src="{{ .Image.Src }}" alt="{{ .Image.Alt }}">
</li>
  {{ end }}
{{ end }}
</ul>
