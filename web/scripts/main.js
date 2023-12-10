"use strict";

export function getCSRFToken() {
  return document.querySelector('meta[name="csrf_"]').content;
}

export function getRegisterData() {
  return {
    showModal: false,
    username: "",
    password: "",
    confirmPassword: "",
    usernameReserved: false,
    notifs: {
      u: false,
      pw: false,
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
    sanitizeUsername,
    isUsernameValid,
    isPasswordValid,
    setStrengths,
    setEvalStrengths,
  };
}

export function getResetPasswordData() {
  return {
    username: "",
    password: "",
    confirmPassword: "",
    notifs: {
      u: true,
      pw: false,
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
    sanitizeUsername,
    isUsernameValid,
    isPasswordValid,
    setStrengths,
    setEvalStrengths,
  };
}

export function sanitizeUsername(u) {
  return u.replace(/[^a-zA-Z0-9_-]+/gi, "").toLowerCase();
}

// TODO: Test
export function sanitizeCharacterName(n) {
  return n.replace(/[^a-zA-Z0-9_-]+/gi, "");
}

// TODO: Pass these lengths in as constants
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

// TODO: Test
export function isCharacterNameValid(n = "") {
  if (n.length < 4) return false;
  if (n.length > 16) return false;
  const regex = new RegExp("[^a-zA-Z0-9_-]+", "g");
  if (regex.test(n)) return false;
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
    showModal: false,
    username: "",
    password: "",
    sanitizeUsername,
  };
}

export function getProfileEmailData() {
  return {
    addEmailMode: false,
    addEmail: "",
  };
}

export function getEmailData(email) {
  return {
    loadEmail: email,
    email,
    editMode: false,
    deleteMode: false,
  };
}

export function getGravatarEmailData(selectedEmail) {
  return {
    selectedEmail,
  };
}

export function getProfileAvatarData(
  avatarSource,
  gravatarHash,
  githubUsername,
) {
  return {
    avatarSource,
    gravatarHash,
    githubUsername,
    getProfileAvatarSrc,
  };
}

export function getProfileAvatarSrc(
  avatarSource,
  gravatarHash,
  githubUsername,
) {
  if (avatarSource === "github") {
    return `https://github.com/${githubUsername}.png`;
  } else {
    return `https://gravatar.com/avatar/${gravatarHash}.jpeg?s=256&d=retro`;
  }
}

export function getCharacterApplicationFlowNameData(name) {
  return {
    name,
    eval: {
      n: {
        len: false,
      },
    },
    sanitizeCharacterName,
    isCharacterNameValid,
  };
}

const HEADER_CSRF_TOKEN = "X-CSRF-Token";
const HEADER_HX_ACCEPTABLE = "X-HX-Acceptable";
const HX_ACCEPTABLE_STATUSES = {
  400: true,
  401: true,
  403: true,
  404: true,
  409: true,
  500: true,
};

document.body.addEventListener("htmx:configRequest", (event) => {
  event.detail.headers[HEADER_CSRF_TOKEN] = getCSRFToken();
});

document.body.addEventListener("htmx:beforeOnLoad", (event) => {
  if (event.detail.xhr.getResponseHeader(HEADER_HX_ACCEPTABLE) !== "true") {
    return;
  }

  if (HX_ACCEPTABLE_STATUSES[event.detail.xhr.status]) {
    event.detail.shouldSwap = true;
    event.detail.isError = false;
  }
});

window.getCSRFToken = getCSRFToken;
window.getRegisterData = getRegisterData;
window.getResetPasswordData = getResetPasswordData;
window.getLoginData = getLoginData;
window.getProfileEmailData = getProfileEmailData;
window.getEmailData = getEmailData;
window.getGravatarEmailData = getGravatarEmailData;
window.getProfileAvatarData = getProfileAvatarData;
window.getCharacterApplicationFlowNameData =
  getCharacterApplicationFlowNameData;
