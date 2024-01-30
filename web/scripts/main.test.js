import { describe, test, expect, beforeEach } from "vitest";
import { isUsernameValid, setStrengths } from "./main.js";

describe("isUsernameValid", () => {
  describe("Invalid if", () => {
    test("username is too short", () => {
      const username = "tes";
      expect(isUsernameValid(username)).toBe(false);
    });
    test("username is too long", () => {
      const username = "testtestt";
      expect(isUsernameValid(username)).toBe(false);
    });
    test("username contains invalid characters", () => {
      const usernameOne = "Test";
      expect(isUsernameValid(usernameOne)).toBe(false);
      const usernameTwo = "test^";
      expect(isUsernameValid(usernameTwo)).toBe(false);
      const usernameThree = "test&*#";
      expect(isUsernameValid(usernameThree)).toBe(false);
    });
  });
  describe("Valid if", () => {
    test("username is the correct length", () => {
      const username = "test";
      expect(isUsernameValid(username)).toBe(true);
    });
    test("username contains numbers", () => {
      const username = "test4u";
      expect(isUsernameValid(username)).toBe(true);
    });
    test("username contains dashes", () => {
      const username = "test-u";
      expect(isUsernameValid(username)).toBe(true);
    });
    test("username contains underscores", () => {
      const username = "test_u";
      expect(isUsernameValid(username)).toBe(true);
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
      expect(strengths.len).toBe(true);
    });
    test("Lacks girth", () => {
      const pw = "shortpw";
      setStrengths(strengths, pw);
      expect(strengths.len).toBe(false);
    });
  });
  describe("Mixed Case", () => {
    test("Strong with mixed case", () => {
      const pw = "tEst";
      setStrengths(strengths, pw);
      expect(strengths.mixedCase).toBe(true);
    });
    test("Weak with all lowercase", () => {
      const pw = "test";
      setStrengths(strengths, pw);
      expect(strengths.mixedCase).toBe(false);
    });
    test("Weak with all uppercase", () => {
      const pw = "TEST";
      setStrengths(strengths, pw);
      expect(strengths.mixedCase).toBe(false);
    });
  });
  describe("Number", () => {
    test("Strong with number", () => {
      const pw = "test1";
      setStrengths(strengths, pw);
      expect(strengths.num).toBe(true);
    });
    test("Strong with numbers", () => {
      const pw = "test12";
      setStrengths(strengths, pw);
      expect(strengths.num).toBe(true);
    });
    test("Weak without numbers", () => {
      const pw = "test";
      setStrengths(strengths, pw);
      expect(strengths.num).toBe(false);
    });
  });
  describe("Special Characters", () => {
    test("Strong with special character", () => {
      const pw = "~";
      setStrengths(strengths, pw);
      expect(strengths.specialChar).toBe(true);
    });
    test("Weak without special character", () => {
      const pw = "test123";
      setStrengths(strengths, pw);
      expect(strengths.specialChar).toBe(false);
    });
  });
});

describe("isActorImageNameValid", () => {
  describe("Invalid if", () => {
    test("Actor Image name is too short", () => {
      const name = "tes";
      expect(isActorImageNameValid(name)).toBe(false);
    });
    test("Actor Image name is too long", () => {
      const name =
        "actor-image-name-of-incredible-proportions-far-exceeding-any-reasonable-character-limit";
      expect(isActorImageNameValid(name)).toBe(false);
    });
    test("Actor Image name contains invalid characters", () => {
      const nameOne = "Test";
      expect(isActorImageNameValid(nameOne)).toBe(false);
      const nameTwo = "test^";
      expect(isActorImageNameValid(nameTwo)).toBe(false);
      const nameThree = "test&*#";
      expect(isActorImageNameValid(nameThree)).toBe(false);
      const nameFour = "test4";
      expect(isActorImageNameValid(nameFour)).toBe(false);
      const nameFive = "test_actor_image";
      expect(isActorImageNameValid(nameFive)).toBe(false);
    });
  });
  describe("Valid if", () => {
    test("Actor Image name is the correct length", () => {
      const name = "test";
      expect(isActorImageNameValid(name)).toBe(true);
    });
    test("Actor Image name contains dashes", () => {
      const name = "test-actor-image";
      expect(isActorImageNameValid(name)).toBe(true);
    });
  });
});
