// Section heading partial: light accent block with dark accent text
#let section-heading(title, theme, to-length) = {
  let bg = rgb(theme.colors.accent_bg)
  let fg = rgb(theme.colors.accent_fg)
  let section-size = to-length(theme.typography.font_size_section)
  let section-gap = to-length(theme.spacing.section_gap)
  let margin-x = 10pt

  v(section-gap * 0.4)
  // Negative padding bleeds the background into page margins;
  // text stays centered in the normal text area.
  pad(left: -margin-x, right: -margin-x,
    block(
      width: 100%,
      fill: bg,
      inset: (x: margin-x, y: 3.5pt),
      align(center, text(
        fill: fg,
        weight: "bold",
        size: section-size,
        tracking: 1pt,
        upper(title),
      )),
    ),
  )
  v(section-gap * 0.2)
}
