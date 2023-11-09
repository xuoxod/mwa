// Verify phone
const verifyPhoneForm = document.querySelector("#verify-phone-form");
const verifyPhoneButton = document.querySelector("#verify-phone-button");
const verifyPhoneInput = document.querySelector("#verify-phone-input");
const verifyPhoneLabel = document.querySelector("#verify-phone-label");
const phoneLabelIcon = document.querySelector("#phone-label-icon");

const sendPhoneVerification = () => {
  log(`Sending phone verification code\n`);
};

const verifyPhoneButtonHandler = () => {
  phoneLabelIcon.remove();

  const p = newElement("p");
  p.innerText = "A verification code has been sent to your phone";

  addAttribute(verifyPhoneInput, "placeholder", "Enter verification code");
  addAttribute(p, "class", "text-success fw-bold m-1");
  addAttribute(p, "id", "p-phone");

  appendChild(verifyPhoneLabel, p);

  verifyPhone();
};

function verifyPhone() {
  log(`Requested phone verification`);

  const url = "/user/phone/verify";

  try {
    fetch(url)
      .then((response) => response.json())
      .then((data) => {
        phoneVerificationResponse(data);
      });
  } catch (err) {
    log(err);
  }
}

function phoneVerificationResponse(response) {
  if (response["ok"]) {
    log(`Phone verification succeeded\n`);

    verifyPhoneInput.classList.remove("d-none");
    document.querySelector("#p-phone").remove();
    phoneIcon = newElement("i");

    addAttribute(
      phoneIcon,
      "class",
      "bi bi-telephone-fill fw-bold fs-3 text-success"
    );
    verifyPhoneButton.classList.remove("btn-danger");
    verifyPhoneButton.classList.add("btn-success");
    verifyPhoneButton.innerText = "Submit";
    verifyPhoneButton.removeEventListener("click", verifyPhoneButtonHandler);
    verifyPhoneButton.addEventListener("click", sendPhoneVerification);

    if (countChildren(verifyPhoneLabel) > 0) {
      removeChildren(verifyPhoneLabel);
    }

    appendChild(verifyPhoneLabel, phoneIcon);
  } else {
    log(`Phone verification failed\n`);
  }
}

// Verify email
const verifyEmailForm = document.querySelector("#verify-email-form");
const verifyEmailButton = document.querySelector("#verify-email-button");
const verifyEmailInput = document.querySelector("#verify-email-input");
const verifyEmailLabel = document.querySelector("#verify-email-label");
const emailLabelIcon = document.querySelector("#email-label-icon");

const sendEmailVerification = () => {
  log(`Sending email verification code\n`);
};

const verifyEmailButtonHandler = () => {
  emailLabelIcon.remove();

  const p = newElement("p");
  p.innerText = "A verification code has been sent to your email";

  addAttribute(verifyEmailInput, "placeholder", "Enter verification code");
  addAttribute(p, "class", "text-success fw-bold m-1");
  addAttribute(p, "id", "p-email");

  appendChild(verifyEmailLabel, p);

  verifyEmail();
};

function verifyEmail() {
  log(`Requested email verification`);

  const url = "/user/email/verify";

  try {
    fetch(url)
      .then((response) => response.json())
      .then((data) => {
        emailVerificationResponse(data);
      });
  } catch (err) {
    log(err);
  }
}

function emailVerificationResponse(response) {
  if (response["ok"]) {
    log(`Email verification succeeded\n`);

    verifyEmailInput.classList.remove("d-none");
    document.querySelector("#p-email").remove();
    emailIcon = newElement("i");

    addAttribute(
      emailIcon,
      "class",
      "bi bi-envelope-at-fill fw-bold fs-3 text-success"
    );
    verifyEmailButton.classList.remove("btn-danger");
    verifyEmailButton.classList.add("btn-success");
    verifyEmailButton.innerText = "Submit";
    verifyEmailButton.removeEventListener("click", verifyEmailButtonHandler);
    verifyEmailButton.addEventListener("click", sendEmailVerification);

    if (countChildren(verifyEmailLabel) > 0) {
      removeChildren(verifyEmailLabel);
    }

    appendChild(verifyEmailLabel, emailIcon);
  } else {
    log(`Email verification failed\n`);
  }
}

// Register click event
verifyPhoneButton.addEventListener("click", verifyPhoneButtonHandler);

verifyEmailButton.addEventListener("click", verifyEmailButtonHandler);
