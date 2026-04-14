// Rizzume — ATS-friendly single-page resume template
//
// Expects two --input arguments:
//   content-path: path to content JSON file
//   theme-path:   path to theme JSON file

#let content = json(sys.inputs.at("content-path"))
#let theme = json(sys.inputs.at("theme-path"))

// Helper: convert dimension strings ("10pt", "0.45in") to Typst lengths
#let to-length(s) = eval(s, mode: "code")

// Import partials
#import "partials/header.typ": render-header
#import "partials/section.typ": section-heading
#import "partials/skills.typ": render-skills
#import "partials/experience.typ": render-experience
#import "partials/education.typ": render-education

// --- PDF metadata (helps ATS and search) ---
#set document(
  title: content.basics.name + " — Resume",
  author: content.basics.name,
  keywords: if "skills" in content {
    content.skills.groups.map(g => g.name).join(", ")
  } else { "" },
)

// --- Page setup ---
#set page(
  width: to-length(theme.page.width),
  height: to-length(theme.page.height),
  margin: (
    x: to-length(theme.page.margin_x),
    y: to-length(theme.page.margin_y),
  ),
)

// --- Typography ---
#set text(
  font: theme.typography.font_body,
  size: to-length(theme.typography.font_size_base),
  fill: rgb(theme.colors.text),
  lang: "en",
)

#set par(
  leading: (theme.typography.line_height - 1.0) * to-length(theme.typography.font_size_base) + 0.65em,
  spacing: 0.65em,
)

// Disable default heading numbering/styling (we use custom section headings)
#set heading(numbering: none)
#show heading: it => it.body

// Links: inherit surrounding text style (no underline or color change)
#show link: it => it

// --- Render resume ---

#render-header(content, theme, to-length)

// Summary
#if "summary" in content {
  section-heading("Summary", theme, to-length)
  text(content.summary.text)
}

// Skills
#if "skills" in content {
  section-heading("Skills", theme, to-length)
  render-skills(content, theme, to-length)
}

// Experience
#if "experience" in content {
  section-heading("Experience", theme, to-length)
  render-experience(content, theme, to-length)
}

// Education
#if "education" in content {
  section-heading("Education", theme, to-length)
  render-education(content, theme, to-length)
}
