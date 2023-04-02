// function get data from contact form
function getData() {
  let name = document.getElementById("name").value
  let email = document.getElementById("email").value
  let phone = document.getElementById("phone").value
  let subject = document.getElementById("subject").value
  let message = document.getElementById("message").value
  
  // validation if form is empty
  if(name == "") {
    return alert("Name must be filled!")
  } else if(email == "") {
    return alert("Email must be filled")
  } else if(phone = "") {
    return alert("Phone Number must be filled")
  } else if(subject == "") {
    return alert("Subject must be filled")
  } else if(message = "") {
    return alert("Message must be filled")
  }
  
  // destination email 
  const destination = "mochammadtaufiq779@gmail.com"

  let a = document.createElement("a")
  a.href = `mailto:${destination}?subject=${subject}&body= Hello my name is ${name}, i want ${message}. You can contact me at this number ${phone} or at email ${email}`

  a.click()

  let data = {
    name: name,
    email: email,
    phone: phone,
    subject: subject,
    message:message
  }
  console.log(data)
}