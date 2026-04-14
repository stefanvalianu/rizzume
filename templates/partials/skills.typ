// Skills partial: grid with bold labels and items
#let render-skills(content, theme, to-length) = {
  let groups = content.skills.groups
  let label-width = to-length(theme.spacing.skills_label_width)
  let item-gap = to-length(theme.spacing.item_gap)

  set text(size: to-length(theme.typography.font_size_base))

  for (i, group) in groups.enumerate() {
    grid(
      columns: (label-width, 1fr),
      gutter: 8pt,
      text(weight: "bold", group.name),
      text(group.items),
    )
    if i < groups.len() - 1 {
      v(item-gap * 0.5)
    }
  }
}
