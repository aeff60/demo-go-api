
fetch('https://5c45-223-24-159-173.ngrok.io/course')
  .then(resp => resp.json())
  .then(data => displayData(data));

function displayData(data) {
  document.querySelector("pre").innerHTML = JSON.stringify(data);
}