"use strict";

export function getRegisterData() {
  return {
    username: "",
    password: "",
    confirmPassword: "",
    pwLen: false,
    pwEvalLen: false,
    pwMixedCase: false,
    pwEvalMixedCase: false,
    pwNum: false,
    pwEvalNum: false,
    pwSpecialChar: false,
    pwEvalSpecialChar: false,
    submitData: async () => {
      if (!isUsernameValid(this.username)) return;
      if (!isPasswordValid(this.password)) return;
      if (this.password !== this.confirmPassword) return;
      const u = sanitizeUsername(this.username);

      const body = new FormData();
      body.append("username", u);
      body.append("password", password);

      const response = await fetch("/player/new", {
        method: "POST",
        body,
      });

      console.dir(response.status);
    },
  };
}

export function sanitizeUsername(u) {
  return u.replace(/[^a-zA-Z0-9_\-]+/gi, "").toLowerCase();
}

export function isUsernameValid(u) {
  if (u.length < 4) return false;
  if (u.length > 16) return false;
  const regex = new RegExp("[^a-z0-9_-]+", "g");
  if (regex.test(u)) return false;
  return true;
}

export function isPasswordValid(pw) {
  if (pw.length < 8) return false;
  if (pw.length > 255) return false;
  return true;
}

export function getPasswordStrengths(pw) {
  let strengths = {
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

  return strengths;
}

window.getRegisterData = getRegisterData;
window.sanitizeUsername = sanitizeUsername;
window.isUsernameValid = isUsernameValid;
window.getPasswordStrengths = getPasswordStrengths;
