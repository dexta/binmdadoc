## 2024-10-09 - Accessible Tooling Buttons
**Learning:** Adding semantic `<button>` tags instead of `<span onclick>` instantly enables keyboard interaction (Space/Enter). The custom 'code-bedder' element requires strict CSS `!important` directives to prevent external themes (like PrismJS) from breaking the textarea-to-code alignment overlay.
**Action:** Always prefer semantic HTML elements over attached click handlers on non-interactive tags. When building overlaying editors, aggressively scope layout/typography variables and use `wrap='off'` on textareas.
