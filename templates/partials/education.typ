// Education partial
#let render-education(content, theme, to-length) = {
  let entries = content.education
  let item-gap = to-length(theme.spacing.item_gap)
  let base-size = to-length(theme.typography.font_size_base)
  let muted = rgb(theme.colors.muted)

  set text(size: base-size)

  for (i, entry) in entries.enumerate() {
    grid(
      columns: (1fr, auto),
      align: (left, right),
      {
        text(weight: "bold", entry.degree)
        [ #sym.bullet ]
        text(entry.school)
        [ #sym.bullet ]
        text(entry.location)
      },
      text(
        style: "italic",
        fill: muted,
        entry.year,
      ),
    )
    if i < entries.len() - 1 {
      v(item-gap)
    }
  }
}
