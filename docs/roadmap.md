Telegram rich text mapping should be done for these:
- spoiler text
- datetime text
- code text
- url text

- custom emoji text
- mathematical expression text
- email address text
- phone number text
- bank card number text
- mention text
- hashtag text
- cashtag text
- bot command text
- button
- anchor
- acnhor link
- reference
- reference link
- text mention text

Platejs node types to do:
- a
- h1-h6
- code (inline)
- code_block
- inline_equation
- blockquote
- table
- th
- tr
- td
- toc

- img
- file
- audio

-- listStyleType

- disc
- -- suggestion

- insert
- remove
- block

- mention

---

Platejs node type 'p' can have nested node types of:
- a
- date
- inline_equation

---

1. Paragraph **'p'**, **'h1'-'h6'**, **'blockquote'**, **'a'** (on its child text), **'li'** (on its child text), **'td', 'th'**, **'code_block'** type can be:
    - bold
    - italic
    - underline
    - strikethrough
    - code (inline)
    - superscript
    - subscript
    - highlight
    - kbd

2. 


Plate.js to Telegram Rich Messages mapping:
- highlight == text marked
- 