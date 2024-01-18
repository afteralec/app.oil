"use strict";

export function toggleTheme() {
  if (document.body.classList.contains("light")) {
    document.body.classList.remove("light");
    document.body.classList.add("dark");
  } else if (document.body.classList.contains("dark")) {
    document.body.classList.remove("dark");
    document.body.classList.add("light");
  }
}

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
  return n.replace(/[^a-zA-Z'-]+/gi, "");
}

// TODO: Test
export function sanitizeCharacterShortDescription(sdesc) {
  return sdesc.replace(/[^a-zA-Z, -]+/gi, "").toLowerCase();
}

// TODO: Test
export function sanitizeCharacterDescription(desc) {
  return desc.replace(/[^a-zA-Z, '-.!()]+/gi, "");
}

// TODO: Test
export function sanitizeCharacterBackstory(bs) {
  return bs.replace(/[^a-zA-Z, '\-\.!()\r\n]+/gi, "");
}

export function sanitizeRequestComment(c = "") {
  const regex = /[^a-zA-Z, "'\-\.?!()\r\n]+/gi;
  return c.replace(regex, "");
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
  const regex = new RegExp("[^a-zA-Z'-]+", "g");
  if (regex.test(n)) return false;
  return true;
}

// TODO: Test
export function isCharacterShortDescriptionValid(sdesc = "") {
  if (sdesc.length < 8) return false;
  if (sdesc.length > 300) return false;
  const regex = new RegExp("[^a-zA-Z, -]+", "g");
  if (regex.test(sdesc)) return false;
  return true;
}

// TODO: Test
export function isCharacterDescriptionValid(desc = "") {
  if (desc.length < 32) return false;
  if (desc.length > 2000) return false;
  const regex = new RegExp("[^a-zA-Z, '-.!()]+", "g");
  if (regex.test(desc)) return false;
  return true;
}

// TODO: Test
export function isCharacterBackstoryValid(bs) {
  if (bs.length < 500) return false;
  if (bs.length > 10000) return false;
  const regex = /[^a-zA-Z, '\-\.!()\r\n]+/gi;
  if (regex.test(bs)) return false;
  return true;
}

// TODO: Test
export function isRequestCommentValid(c = "") {
  if (c.length < 1) return false;
  if (c.length > 500) return false;
  const regex = /[^a-zA-Z, "'\-\.?!()\r\n]+/gi;
  if (regex.test(c)) return false;
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

export function getProfileChangePasswordData() {
  return {
    current: "",
    password: "",
    confirmPassword: "",
    notifs: {
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
    isPasswordValid,
    setStrengths,
    setEvalStrengths,
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
    return `https://gravatar.com/avatar/${gravatarHash}.jpeg?f=y&r=m&s=256&d=retro`;
  }
}

export function parseTimeStamp(ts) {
  const d = new Date(ts * 1000).toLocaleString("en-us", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "numeric",
  });
  return d.toString();
}

export function getCharacterApplicationNameData(name) {
  return {
    name,
    eval: {
      n: {
        len: name.length > 0,
      },
    },
    sanitizeCharacterName,
    isCharacterNameValid,
  };
}

export function getCharacterApplicationGenderData(gender) {
  return {
    gender,
    showSubmitDialog: false,
    showCancelDialog: false,
    partsOpen: false,
    actionsOpen: false,
    commentOpen: false,
  };
}

export function getCharacterApplicationShortDescriptionData(sdesc) {
  return {
    sdesc,
    eval: {
      sdesc: {
        len: sdesc.length > 0,
      },
    },
    sanitizeCharacterShortDescription,
    isCharacterShortDescriptionValid,
    showSubmitDialog: false,
    showCancelDialog: false,
    partsOpen: false,
    actionsOpen: false,
    commentOpen: false,
  };
}

export function getCharacterApplicationDescriptionData(description) {
  return {
    description,
    eval: {
      desc: {
        len: description.length > 0,
      },
    },
    sanitizeCharacterDescription,
    isCharacterDescriptionValid,
    showSubmitDialog: false,
    showCancelDialog: false,
    partsOpen: false,
    actionsOpen: false,
    commentOpen: false,
  };
}

export function getCharacterApplicationBackstoryData(backstory) {
  return {
    backstory,
    eval: {
      bs: {
        len: backstory.length > 0,
      },
    },
    sanitizeCharacterBackstory,
    isCharacterBackstoryValid,
    showSubmitDialog: false,
    showCancelDialog: false,
    partsOpen: false,
    actionsOpen: false,
    commentOpen: false,
  };
}

// TODO: Roll this and the generalized func into one?
export function getCharacterApplicationSummaryData() {
  return {
    showSubmitDialog: false,
    showCancelDialog: false,
    partsOpen: true,
    actionsOpen: true,
  };
}

export function getRequestData() {
  return {
    comment: "",
    sanitizeRequestComment,
    isRequestCommentValid,
    showSubmitDialog: false,
    showCancelDialog: false,
    showPutInReviewDialog: false,
    partsOpen: false,
    actionsOpen: false,
    commentOpen: false,
  };
}

export function getCharacterApplicationData() {
  return {
    comment: "",
    sanitizeRequestComment,
    isRequestCommentValid,
    showSubmitDialog: false,
    showCancelDialog: false,
    partsOpen: false,
    actionsOpen: false,
    commentOpen: false,
  };
}

export function getRequestCommentData(c = "") {
  return {
    comment: c,
    editOpen: false,
    showDeleteConfirm: false,
  };
}

export function getPlayerPermissionsData() {
  return {
    username: "",
    sanitizeUsername,
  };
}

export function getSearchHelpIndexData() {
  return {
    search: "",
    tags: false,
    category: false,
    title: true,
    content: true,
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

window.toggleTheme = toggleTheme;
window.getCSRFToken = getCSRFToken;
window.getRegisterData = getRegisterData;
window.getResetPasswordData = getResetPasswordData;
window.getLoginData = getLoginData;
window.getProfileEmailData = getProfileEmailData;
window.getEmailData = getEmailData;
window.getProfileChangePasswordData = getProfileChangePasswordData;
window.getGravatarEmailData = getGravatarEmailData;
window.getProfileAvatarData = getProfileAvatarData;
window.parseTimeStamp = parseTimeStamp;
window.getCharacterApplicationNameData = getCharacterApplicationNameData;
window.getCharacterApplicationGenderData = getCharacterApplicationGenderData;
window.getCharacterApplicationShortDescriptionData =
  getCharacterApplicationShortDescriptionData;
window.getCharacterApplicationDescriptionData =
  getCharacterApplicationDescriptionData;
window.getCharacterApplicationBackstoryData =
  getCharacterApplicationBackstoryData;
// TODO: I believe this function can be removed
window.getCharacterApplicationSummaryData = getCharacterApplicationSummaryData;
// TODO: I think this function can be removed
window.getCharacterApplicationData = getCharacterApplicationData;
window.getRequestData = getRequestData;
window.getPlayerPermissionsData = getPlayerPermissionsData;
window.getRequestCommentData = getRequestCommentData;
window.getSearchHelpIndexData = getSearchHelpIndexData;
