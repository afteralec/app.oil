function sanitizeUsername(s){return s.replace(/[^a-zA-Z0-9_\-]+/gi,"").toLowerCase()}function isUsernameValid(s){if(s.length<4)return!1;if(s.length>8)return!1;if(new RegExp("[^a-z0-9_-]+","g").test(s))return!1;return!0}function getPasswordStrengths(pw){const strengths={len:!1,mixedCase:!1,num:!1,specialChar:!1};if(pw.length>8)strengths.len=!0;if(pw.match(/[a-z]/)&&pw.match(/[A-Z]/))strengths.mixedCase=!0;if(pw.match(/[0-9]/))strengths.num=!0;if(pw.match(/[^a-zA-Z\d]/))strengths.specialChar=!0;return console.dir(strengths),strengths}export{sanitizeUsername,isUsernameValid,getPasswordStrengths};
