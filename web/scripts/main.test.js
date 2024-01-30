import { describe, test, expect, beforeEach } from "bun:test";
import { isUsernameValid, setStrengths } from "./main.js";

describe("isUsernameValid", () => {
  describe("Invalid if", () => {
    test("username is too short", () => {
      const username = "tes";
      expect(isUsernameValid(username)).toBeFalse();
    });
    test("username is too long", () => {
      const username = "testtestt";
      expect(isUsernameValid(username)).toBeFalse();
    });
    test("username contains invalid characters", () => {
      const usernameOne = "Test";
      expect(isUsernameValid(usernameOne)).toBeFalse();
      const usernameTwo = "test^";
      expect(isUsernameValid(usernameTwo)).toBeFalse();
      const usernameThree = "test&*#";
      expect(isUsernameValid(usernameThree)).toBeFalse();
    });
  });
  describe("Valid if", () => {
    test("username is the correct length", () => {
      const username = "test";
      expect(isUsernameValid(username)).toBeTrue();
    });
    test("username contains numbers", () => {
      const username = "test4u";
      expect(isUsernameValid(username)).toBeTrue();
    });
    test("username contains dashes", () => {
      const username = "test-u";
      expect(isUsernameValid(username)).toBeTrue();
    });
    test("username contains underscores", () => {
      const username = "test_u";
      expect(isUsernameValid(username)).toBeTrue();
    });
  });
});

describe("setStrengths", () => {
  let strengths;
  beforeEach(() => {
    strengths = {
      len: false,
      mixedCase: false,
      num: false,
      specialChar: false,
    };
  });
  describe("Length", () => {
    test("Long is strong", () => {
      const pw = "longpassword";
      setStrengths(strengths, pw);
      expect(strengths.len).toBeTrue();
    });
    test("Lacks girth", () => {
      const pw = "shortpw";
      setStrengths(strengths, pw);
      expect(strengths.len).toBeFalse();
    });
  });
  describe("Mixed Case", () => {
    test("Strong with mixed case", () => {
      const pw = "tEst";
      setStrengths(strengths, pw);
      expect(strengths.mixedCase).toBeTrue();
    });
    test("Weak with all lowercase", () => {
      const pw = "test";
      setStrengths(strengths, pw);
      expect(strengths.mixedCase).toBeFalse();
    });
    test("Weak with all uppercase", () => {
      const pw = "TEST";
      setStrengths(strengths, pw);
      expect(strengths.mixedCase).toBeFalse();
    });
  });
  describe("Number", () => {
    test("Strong with number", () => {
      const pw = "test1";
      setStrengths(strengths, pw);
      expect(strengths.num).toBeTrue();
    });
    test("Strong with numbers", () => {
      const pw = "test12";
      setStrengths(strengths, pw);
      expect(strengths.num).toBeTrue();
    });
    test("Weak without numbers", () => {
      const pw = "test";
      setStrengths(strengths, pw);
      expect(strengths.num).toBeFalse();
    });
  });
  describe("Special Characters", () => {
    test("Strong with special character", () => {
      const pw = "~";
      setStrengths(strengths, pw);
      expect(strengths.specialChar).toBeTrue();
    });
    test("Weak without special character", () => {
      const pw = "test123";
      setStrengths(strengths, pw);
      expect(strengths.specialChar).toBeFalse();
    });
  });
});

describe("isActorImageNameValid", () => {
  describe("Invalid if", () => {
    test("Actor Image name is too short", () => {
      const name = "tes";
      expect(isActorImageNameValid(name)).toBeFalse();
    });
    test("Actor Image name is too long", () => {
      const name =
        "actor-image-name-of-incredible-proportions-far-exceeding-any-reasonable-character-limit";
      expect(isActorImageNameValid(name)).toBeFalse();
    });
    test("Actor Image name contains invalid characters", () => {
      const nameOne = "Test";
      expect(isActorImageNameValid(nameOne)).toBeFalse();
      const nameTwo = "test^";
      expect(isActorImageNameValid(nameTwo)).toBeFalse();
      const nameThree = "test&*#";
      expect(isActorImageNameValid(nameThree)).toBeFalse();
      const nameFour = "test4";
      expect(isActorImageNameValid(nameFour)).toBeFalse();
      const nameFive = "test_actor_image";
      expect(isActorImageNameValid(nameFive)).toBeFalse();
    });
  });
  describe("Valid if", () => {
    test("Actor Image name is the correct length", () => {
      const name = "test";
      expect(isActorImageNameValid(name)).toBeTrue();
    });
    test("Actor Image name contains dashes", () => {
      const name = "test-actor-image";
      expect(isActorImageNameValid(name)).toBeTrue();
    });
  });
});
