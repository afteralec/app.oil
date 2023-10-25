export function sanitizeUsername(s) {
  return s.replace(/[^a-zA-Z0-9_\-]+/gi, "").toLowerCase();
}

export function isUsernameValid(s) {
  if (s.length < 4) return false;

  if (s.length > 8) return false;

  const regex = new RegExp("[^a-z0-9_-]+", "g");
  if (regex.test(s)) return false;

  return true;
}

export function getPasswordStrengths(pw) {
  const strengths = {
    len: false,
    mixedCase: false,
    num: false,
    specialChar: false,
  };

  if (pw.length > 8) {
    strengths.len = true;
  }

  if (pw.match(/[a-z]/) && pw.match(/[A-Z]/)) {
    strengths.mixedCase = true;
  }

  if (pw.match(/[0-9]/)) {
    strengths.num = true;
  }

  if (pw.match(/[^a-zA-Z\d]/)) {
    strengths.specialChar = true;
  }

  console.dir(strengths);
  return strengths;
}
