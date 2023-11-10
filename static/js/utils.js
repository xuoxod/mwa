const log = (msg) => console.log(msg);

const stringify = (obj) => JSON.stringify(obj);

const parser = (strObj) => JSON.parse(strObj);

const addHandler = (theElement, whichEvent, method) => {
  if (null != theElement && null != whichEvent && typeof method == "function") {
    theElement.addEventListener(whichEvent, method);
  }
};

const addClickHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "click", handler);
  }
};

const addKeyupHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "keyup", handler);
  }
};

const addKeydownHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "keydown", handler);
  }
};

const addOnFocusHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "focus", handler);
  }
};

const addOnChangeHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "change", handler);
  }
};

const addOffFocusHandler = (theElement, handler) => {
  if (null != theElement && typeof handler == "function") {
    addHandler(theElement, "focusout", handler);
  }
};

const appendChild = (parent, child) => {
  if (null != parent && null != child) {
    parent.appendChild(child);
  }
};

const appendBeforeLastChild = (parent, child) => {
  if (null != parent && null != child) {
    const lastChildIndex = parent.children.length - 1;
    const lastChild = parent.children[lastChildIndex];
    parent.insertBefore(child, lastChild);
  }
};

const append = (parent, child) => {
  parent.append(child);
};

const removeChildren = (parent) => {
  parent.querySelectorAll("*").forEach((dialog) => {
    dialog.remove();
  });
};

const countChildren = (parent) => {
  if (null != parent) {
    return parent.children.length;
  }
  return null;
};

const addAttribute = (theElement, whichAttribute, attributeValue) => {
  if (null != theElement) {
    theElement.setAttribute(whichAttribute, attributeValue);
  }
};

const setAttribute = (theElement, whichAttribute, attributeValue) => {
  if (null != theElement) {
    theElement.setAttribute(whichAttribute, attributeValue);
  }
};

const getAttribute = (theElement, whichAttribute) => {
  if (null != theElement && null != whichAttribute) {
    return theElement.getAttribute(`${whichAttribute}`) || null;
  }
  return "Element is null";
};

const removeAttribute = (theElement, whichAttribute) => {
  if (null != theElement) {
    if (theElement.hasAttribute(whichAttribute)) {
      theElement.removeAttribute(whichAttribute);
    }
  }
};

const getElement = (nameIdClass) => {
  let element = null;
  if (null != (element = document.querySelector(`${nameIdClass}`))) {
    return element;
  }
  if (null != (element = document.querySelector(`#${nameIdClass}`))) {
    return element;
  }
  if (null != (element = document.querySelector(`.${nameIdClass}`))) {
    return element;
  }
  return null;
};

const cap = (arg) => {
  let word_split = null,
    line = "",
    word = arg.toString();
  if (null !== word && undefined !== word) {
    if (
      word.trim().toLowerCase() === "id" ||
      word.trim().toLowerCase() === "ssn" ||
      word.trim().toLowerCase() === "sku" ||
      word.trim().toLowerCase() === "vm" ||
      word.trim().toLowerCase() === "mac" ||
      word.trim().toLowerCase() === "imei" ||
      word.trim().toLowerCase() === "os" ||
      word.trim().toLowerCase() === "atm" ||
      word.trim().toLowerCase() === "pa" ||
      word.trim().toLowerCase() === "rjw"
    ) {
      return word.toUpperCase();
    } else if (word.match(/[-]/)) {
      if (null !== (word_split = word.split(["-"])).length > 0) {
        for (let i = 0; i < word_split.length; i++) {
          if (i < word_split.length - 1) {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1) +
              "-";
          } else {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1);
          }
        }
        return line;
      }
    } else if (word.match(/[ ]/)) {
      if (null !== (word_split = word.split([" "])).length > 0) {
        for (let i = 0; i < word_split.length; i++) {
          if (i < word_split.length - 1) {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1) +
              " ";
          } else {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1);
          }
        }
        return line;
      }
    } else {
      return word.substring(0, 1).toUpperCase() + word.substring(1);
    }
  }
};

const newElement = (type) => {
  if (null != type && typeof type == "string") {
    return document.createElement(type);
  }
  return null;
};

// notify display a custom alert to user
// type string: success, error or warning
// msg string: message to user
const notify = (type, msg, time = 3) => {
  notie.alert({
    type: type,
    text: msg,
    time,
  });
};

// modal display custom modal
// title string: modal's title
// text string: the message to user
// icon built-in: success, warning, error, info or question
// btnText string: button's text
const modal = (
  title,
  text,
  icon,
  btnText = "Okay",
  btnClose = true,
  showStatus = true
) => {
  Swal.fire({
    title: title,
    text: text,
    icon: icon,
    confirmButtonText: btnText,
    showCloseButton: btnClose,
    showLoading: showStatus,
  });
};
