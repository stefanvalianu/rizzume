// Experience partial: company headers, roles, bullets
#let render-experience(content, theme, to-length) = {
  let entries = content.experience
  let item-gap = to-length(theme.spacing.item_gap)
  let bullet-gap = to-length(theme.spacing.bullet_gap)
  let muted = rgb(theme.colors.muted)
  let base-size = to-length(theme.typography.font_size_base)

  set text(size: base-size)

  for (i, entry) in entries.enumerate() {
    // Company header line
    {
      let company-name = text(weight: "bold", fill: rgb(theme.colors.accent), entry.company)
      let has-url = "url" in entry and entry.url != none and entry.url != ""
      grid(
        columns: (1fr, auto),
        align: (left, right),
        {
          if has-url { link(entry.url, company-name) } else { company-name }
          [ #sym.bullet ]
          text(entry.location)
        },
        text(
          weight: "bold",
          entry.start + [ --- ] + entry.end,
        ),
      )
    }

    let roles = entry.at("roles", default: ())

    for (j, role) in roles.enumerate() {
      let has-bullets = "bullets" in role and role.bullets.len() > 0
      let has-team = "team" in role and role.team != none
      let has-own-dates = "start" in role

      v(item-gap * 0.5)

      // Role line
      pad(left: 10pt, {
        grid(
          columns: (1fr, auto),
          align: (left, right),
          {
            text(weight: "bold", role.title)
            if has-team {
              [ #sym.bullet ]
              text(style: "italic", role.team)
            }
          },
          if has-own-dates {
            text(
              style: "italic",
              fill: muted,
              role.start + [ --- ] + role.end,
            )
          },
        )
      })

      // Bullets
      if has-bullets {
        v(bullet-gap)
        pad(left: 14pt, {
          for (k, bullet) in role.bullets.enumerate() {
            grid(
              columns: (12pt, 1fr),
              align: (left + top, left + top),
              text(sym.bullet),
              text(bullet),
            )
            if k < role.bullets.len() - 1 {
              v(bullet-gap)
            }
          }
        })
      }
    }

    if i < entries.len() - 1 {
      v(item-gap)
    }
  }
}
