async function copyCode(button) {
  const copyButtonLabel = "Copy";
  let figure = button.parentElement.parentElement;
  let text = figure.querySelector("code").innerText;

  await navigator.clipboard.writeText(text);

  // visual feedback that task is completed
  button.innerText = "Copied";
  button.dataset.active = "";
  button.setAttribute("disabled", "");

  setTimeout(() => {
    button.innerText = copyButtonLabel;
    button.removeAttribute("disabled");
    button.removeAttribute("data-active");
  }, 2000);
}

let buttons = document.querySelectorAll("button.copycode");

buttons.forEach((button) => {
  button.addEventListener("click", () => copyCode(button));
  button.removeAttribute("disabled");
});
