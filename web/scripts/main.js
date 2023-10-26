"use strict";

export function getRegisterData() {
  return {
    username: "",
    password: "",
    confirmPassword: "",
    notifs: {
      u: false,
      pw: false,
      cpw: false,
    },
    eval: {
      pw: {
        strengths: {
          len: false,
          mixedCase: false,
          num: false,
          specialChar: false,
        },
      },
      u: {
        len: false,
      },
      cpw: {
        eq: false,
      },
    },
    strengths: {
      len: false,
      mixedCase: false,
      num: false,
      specialChar: false,
    },
    errors: {
      conflict: false,
      internal: false,
      disaster: false,
    },
    submitData,
    sanitizeUsername,
    isUsernameValid,
    isPasswordValid,
    setStrengths,
    setEvalStrengths,
  };
}

export async function submitData(errors, u, pw, confirmpw) {
  if (!isUsernameValid(u)) return;
  if (!isPasswordValid(pw)) return;
  if (pw !== confirmpw) return;
  const su = sanitizeUsername(u);

  const body = new FormData();
  body.append("username", su);
  body.append("password", pw);

  try {
    const response = await fetch("/player/new", {
      method: "POST",
      body,
    });

    if (response.status !== 201) {
      if (response.status === 409) {
        errors.conflict = true;
        return;
      }

      errors.internal = true;
      return;
    }

    window.location.reload();
  } catch (err) {
    errors.disaster = true;
    return;
  }
}

export function sanitizeUsername(u) {
  return u.replace(/[^a-zA-Z0-9_-]+/gi, "").toLowerCase();
}

export function isUsernameValid(u) {
  if (u.length < 4) return false;
  if (u.length > 8) return false;
  const regex = new RegExp("[^a-z0-9_-]+", "g");
  if (regex.test(u)) return false;
  return true;
}

export function isPasswordValid(pw) {
  if (pw.length < 8) return false;
  if (pw.length > 255) return false;
  return true;
}

export function setStrengths(strengths, pw) {
  strengths.len = false;
  if (pw.length > 8) {
    strengths.len = true;
  }

  strengths.mixedCase = false;
  if (pw.match(/[a-z]/) && pw.match(/[A-Z]/)) {
    strengths.mixedCase = true;
  }

  strengths.num = false;
  if (pw.match(/[0-9]/)) {
    strengths.num = true;
  }

  strengths.specialChar = false;
  if (pw.match(/[^a-zA-Z\d]/)) {
    strengths.specialChar = true;
  }
}

export function setEvalStrengths(evalStrengths, strengths) {
  evalStrengths.len = strengths.len || evalStrengths.len;
  evalStrengths.mixedCase = strengths.mixedCase || evalStrengths.mixedCase;
  evalStrengths.num = strengths.num || evalStrengths.num;
  evalStrengths.specialChar =
    strengths.specialChar || evalStrengths.specialChar;
}

export function getLoginData() {
  return {
    username: "",
    password: "",
    errors: {
      auth: false,
      internal: false,
      disaster: false,
    },
    submitLogin,
    sanitizeUsername,
  };
}

export async function submitLogin(errors, u, pw) {
  try {
    const body = new FormData();
    body.append("username", u);
    body.append("password", pw);
    const res = await fetch("/login", {
      method: "POST",
      body,
    });

    if (res.status != 200) {
      if (res.status === 500) {
        errors.internal = true;
        return;
      }

      errors.auth = true;
      return;
    }

    window.location.reload();
  } catch {
    errors.disaster = true;
  }
}

export function getLogoutData() {
  return {
    logout,
  };
}

export async function logout() {
  try {
    const res = await fetch("/logout", {
      method: "POST",
    });

    if (res.status !== 200) {
      // TODO: Something went wrong with destroying the session
      return;
    }

    window.location.reload();
  } catch {
    // TODO: Handle this error case here - backend is unreachable
  }
}

window.getRegisterData = getRegisterData;
window.getLoginData = getLoginData;
window.getLogoutData = getLogoutData;
