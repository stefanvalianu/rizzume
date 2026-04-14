// Header partial: name (split color), contact info with links
#let render-header(content, theme, to-length) = {
  let basics = content.at("basics")
  let accent = rgb(theme.colors.accent)
  let muted = rgb(theme.colors.muted)
  let name-size = to-length(theme.typography.font_size_name)
  let small-size = to-length(theme.typography.font_size_small)

  // Name — last word in accent color
  {
    let name-parts = basics.name.split(" ")
    let first = name-parts.slice(0, -1).join(" ")
    let last = name-parts.last()
    text(size: name-size, weight: "bold", tracking: 0.5pt, upper(first))
    text(size: name-size, weight: "bold", tracking: 0.5pt, [ ])
    text(size: name-size, weight: "bold", tracking: 0.5pt, fill: accent, upper(last))
  }
  linebreak()

  // Contact info line with links
  {
    set text(size: small-size, fill: muted)
    let parts = ()
    if "location" in basics { parts.push(basics.location) }
    if "phone" in basics { parts.push(basics.phone) }
    if "email" in basics {
      parts.push(link("mailto:" + basics.email, basics.email))
    }
    if "linkedin" in basics {
      let url = if basics.linkedin.starts-with("http") { basics.linkedin } else { "https://" + basics.linkedin }
      parts.push(link(url, basics.linkedin))
    }
    if "website" in basics {
      let url = if basics.website.starts-with("http") { basics.website } else { "https://" + basics.website }
      parts.push(link(url, basics.website))
    }
    if "github" in basics {
      let url = if basics.github.starts-with("http") { basics.github } else { "https://" + basics.github }
      parts.push(link(url, basics.github))
    }
    parts.join([ #sym.bullet ])
  }
  v(to-length(theme.spacing.header_gap))
}
