function sanitizeUsername(s) {
  console.log(s);
  return s.replace(/[^a-zA-Z0-9_\-]+/gi, "").toLowerCase();
}
