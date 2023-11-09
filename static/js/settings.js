// Profile tab
const profileSubmitButtonGroup = document.querySelector(
  "#profile-submit-button-group"
);
const profileEditButtonGroup = document.querySelector(
  "#profile-edit-button-group"
);
const profileSubmitButton = document.querySelector("#profile-submit-button");
const profileEditButton = document.querySelector("#profile-edit-button");
const profileForm = document.querySelector("#profile-form");
const inputs = document.querySelectorAll(".form-control");
const fnameGroup = document.querySelector("#fname-group");
const lnameGroup = document.querySelector("#lname-group");
const unameGroup = document.querySelector("#uname-group");
const addressGroup = document.querySelector("#address");
const zipcodeGroup = document.querySelector("#zipcode");
const emailGroup = document.querySelector("#email");
const phoneGroup = document.querySelector("#phone");
const stateGroup = document.querySelector("#state");
const cityGroup = document.querySelector("#city");
let toggler = false;
const objInputs = {};
const changedInputs = {};

inputs.forEach((input) => {
  objInputs[`${input.name.trim()}`] = {
    name: input.name.trim(),
    value: input.value.trim(),
  };
});

const inputValueChanged = (inputName, value) => {
  const input = objInputs[inputName] || null;

  if (null != input) {
    return input.value != value;
  }
  return false;
};

const addChangedInput = (input) => {
  changedInputs[`${input.name}`] = input.value;
};

const removeChangedInput = (inputName) => {
  if (changedInputs.hasOwnProperty(inputName)) {
    delete changedInputs[`${inputName}`];
  }
};

const handleFormError = (data) => {
  if (data["fname"]) {
    const parentCount = countChildren(document.querySelector("#fname-group"));

    if (parentCount > 2) {
      document.querySelector("#fname-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "fname-error");
    span.innerText = data["fname"];
    appendChild(document.querySelector("#fname-group"), span);
  } else {
    if (document.querySelector("#fname-error")) {
      document.querySelector("#fname-error").remove();
    }
  }

  if (data["lname"]) {
    const parentCount = countChildren(document.querySelector("#lname-group"));

    if (parentCount > 2) {
      document.querySelector("#lname-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "lname-error");
    span.innerText = data["lname"];
    appendChild(document.querySelector("#lname-group"), span);
  } else {
    if (document.querySelector("#lname-error")) {
      document.querySelector("#lname-error").remove();
    }
  }

  if (data["uname"]) {
    const parentCount = countChildren(document.querySelector("#uname-group"));

    if (parentCount > 2) {
      document.querySelector("#uname-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "uname-error");
    span.innerText = data["uname"];
    appendChild(document.querySelector("#uname-group"), span);
  } else {
    if (document.querySelector("#uname-error")) {
      document.querySelector("#uname-error").remove();
    }
  }

  if (data["iurl"]) {
    const parentCount = countChildren(document.querySelector("#iurl-group"));

    if (parentCount > 2) {
      document.querySelector("#iurl-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "iurl-error");
    span.innerText = data["iurl"];
    appendChild(document.querySelector("#iurl-group"), span);
  } else {
    if (document.querySelector("#iurl-error")) {
      document.querySelector("#iurl-error").remove();
    }
  }

  if (data["email"]) {
    const parentCount = countChildren(document.querySelector("#email-group"));

    if (parentCount > 2) {
      document.querySelector("#email-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "email-error");
    span.innerText = data["email"];
    appendChild(document.querySelector("#email-group"), span);
  } else {
    if (document.querySelector("#email-error")) {
      document.querySelector("#email-error").remove();
    }
  }

  if (data["phone"]) {
    const parentCount = countChildren(document.querySelector("#phone-group"));

    if (parentCount > 2) {
      document.querySelector("#phone-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "phone-error");
    span.innerText = data["phone"];
    appendChild(document.querySelector("#phone-group"), span);
  } else {
    if (document.querySelector("#phone-error")) {
      document.querySelector("#phone-error").remove();
    }
  }

  if (data["address"]) {
    const parentCount = countChildren(document.querySelector("#address-group"));

    if (parentCount > 2) {
      document.querySelector("#address-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "address-error");
    span.innerText = data["address"];
    appendChild(document.querySelector("#address-group"), span);
  } else {
    if (document.querySelector("#address-error")) {
      document.querySelector("#address-error").remove();
    }
  }

  if (data["city"]) {
    const parentCount = countChildren(document.querySelector("#city-group"));

    if (parentCount > 2) {
      document.querySelector("#city-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "city-error");
    span.innerText = data["city"];
    appendChild(document.querySelector("#city-group"), span);
  } else {
    if (document.querySelector("#city-error")) {
      document.querySelector("#city-error").remove();
    }
  }

  if (data["state"]) {
    const parentCount = countChildren(document.querySelector("#state-group"));

    if (parentCount > 2) {
      document.querySelector("#state-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "state-error");
    span.innerText = data["state"];
    appendChild(document.querySelector("#state-group"), span);
  } else {
    if (document.querySelector("#state-error")) {
      document.querySelector("#state-error").remove();
    }
  }

  if (data["zipcode"]) {
    const parentCount = countChildren(document.querySelector("#zipcode-group"));

    if (parentCount > 2) {
      document.querySelector("#zipcode-group").lastChild.remove();
    }

    const span = newElement("span");
    addAttribute(span, "class", "input-group-text text-danger");
    addAttribute(span, "id", "zipcode-error");
    span.innerText = data["zipcode"];
    appendChild(document.querySelector("#zipcode-group"), span);
  } else {
    if (document.querySelector("#zipcode-error")) {
      document.querySelector("#zipcode-error").remove();
    }
  }
};

const handleFormSuccess = () => {
  log(`User updated\n`);

  // remove all error messages
  if (document.querySelector("#fname-error")) {
    document.querySelector("#fname-error").remove();
  }

  if (document.querySelector("#lname-error")) {
    document.querySelector("#lname-error").remove();
  }

  if (document.querySelector("#uname-error")) {
    document.querySelector("#uname-error").remove();
  }

  if (document.querySelector("#email-error")) {
    document.querySelector("#email-error").remove();
  }

  if (document.querySelector("#phone-error")) {
    document.querySelector("#phone-error").remove();
  }

  if (document.querySelector("#iurl-error")) {
    document.querySelector("#iurl-error").remove();
  }

  if (document.querySelector("#address-error")) {
    document.querySelector("#address-error").remove();
  }

  if (document.querySelector("#city-error")) {
    document.querySelector("#city-error").remove();
  }

  if (document.querySelector("#state-error")) {
    document.querySelector("#state-error").remove();
  }

  if (document.querySelector("#zipcode-error")) {
    document.querySelector("#zipcode-error").remove();
  }

  location.href = "/user/settings";
};

const submitProfileForm = () => {
  const url = "/user/profile";
  const formData = new FormData(profileForm);
  formData.append("csrf_token", "{{.csrftoken}}");

  log(`Sending form data: ${stringify(formData)}\n`);

  try {
    fetch("/user/profile", {
      method: "post",
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        if (data["ok"] == false) {
          handleFormError(data);
        } else {
          handleFormSuccess();
        }
      });
  } catch (err) {
    log(err);
  }
};

profileForm.addEventListener("submit", (e) => {
  e.preventDefault();
  log(`Profile form submitted\n`);
  try {
    submitProfileForm();
  } catch (err) {
    log(err);
  }
  return;
});

profileEditButton.addEventListener("click", () => {
  toggler = !toggler;
  log(`Edit button clicked\n`);
  if (toggler) {
    inputs.forEach((input) => {
      // log(`${input.name}\n`);
      input.removeAttribute("readonly");
    });
    profileEditButton.classList.remove("btn-success");
    profileEditButton.classList.add("btn-warning");
  } else {
    inputs.forEach((input) => {
      // log(`${input.name}\n`);
      if (!input.hasAttribute("readonly")) {
        input.setAttribute("readonly", "");
      }
    });
    profileEditButton.classList.remove("btn-warning");
    profileEditButton.classList.add("btn-success");
  }
});

inputs.forEach((input) => {
  input.addEventListener("keyup", (e) => {
    if (inputValueChanged(e.target.name.trim(), e.target.value.trim())) {
      // log(`Element ${e.target.name}'s value has changed\n`);
      addChangedInput(e.target);
    } else {
      // log(`Element ${e.target.name}'s value was changed back\n`);
      removeChangedInput(e.target.name.trim());
    }

    if (Object.keys(changedInputs).length > 0) {
      profileSubmitButtonGroup.classList.remove("d-none");
      profileEditButton.innerText = "Save";
    } else {
      if (!profileEditButtonGroup.classList.contains("d-none")) {
        profileSubmitButtonGroup.classList.add("d-none");
        profileEditButton.innerText = "Edit";
      }
    }
  });
});

// Settings tab
const settingsForm = document.querySelector("#settings-form");
const handleSettingsFormError = () => {
  log(`Preferences form submission failed\n`);
};
const handleSettingsFormSuccess = () => {
  log(`Preferences form submission succeeded\n`);
  location.href = "/user/settings";
};
const submitSettingsForm = () => {
  const url = "/user/settings";
  const formData = new FormData(settingsForm);
  formData.append("csrf_token", "{{.csrftoken}}");
  log(`Sending form data: ${stringify(formData)}\n`);
  try {
    fetch("/user/settings", {
      method: "post",
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        if (data["ok"] == false) {
          handleSettingsFormError(data);
        } else {
          handleSettingsFormSuccess();
        }
      });
  } catch (err) {
    log(err);
  }
};

settingsForm.addEventListener("submit", (e) => {
  e.preventDefault();
  log(`Settings form submitted\n`);
  try {
    submitSettingsForm();
  } catch (err) {
    log(err);
  }
  return;
});
