// save the result object into an empty array
let blogs = []

function getBlog(event){
  // handle event it so doesn't refresh form onsubmit
  event.preventDefault()

  // handle input from form blog html
  let projectName = document.getElementById("projectName").value
  let description = document.getElementById("description").value
  let startDate = document.getElementById("start-date").value
  let endDate = document.getElementById("end-date").value
  let nodeJs = document.getElementById("nodeJs").checked
  let reactJs = document.getElementById("reactJs").checked
  let nextJs = document.getElementById("nextJs").checked
  let typescript = document.getElementById("typescript").checked
  let image = document.getElementById("inputBlog").files

  // maybe to handle object biner/blob so the image can be displayed
  image = URL.createObjectURL(image[0])
  // Handle result from form, and include in the object
  let blog = {
    projectName: projectName,
    description: description,
    startDate: startDate,
    endDate: endDate,
    nodeJs: nodeJs,
    reactJs: reactJs,
    nextJs: nextJs,
    typescript: typescript,
    image: image
  }
  // inserts objects into an empty array
  blogs.push(blog)
  // make sure the output can be handlle
  console.log(blogs)
  // call function renderBlog()
  renderBlog()
}

function renderBlog() {
  document.getElementById("contents").innerHTML = ""
  // make looping the result of the function getBlog()
  for(let i = 0; i < blogs.length; i++) {
    document.getElementById("contents").innerHTML += `
        <!-- item grid -->
        <div class="grid-item">
          <img class="img-project" src=${blogs[i].image} alt="project">
          <div class="info-project">
            <h3>${blogs[i].projectName}</h3>
            <p style="opacity: 0.8;">${getDuration(blogs[i].startDate, blogs[i].endDate)}</p>
            <p>${blogs[i].description}</p>
          </div>
          <div class="framework-logo">
            ${blogs[i].nextJs ? "<img src='assets/img/nextJs.png'>" : "" }
            ${blogs[i].reactJs ? "<img src='assets/img/reactjs.png'>" : ""}
            ${blogs[i].nodeJs ? "<img src='assets/img/nodejs.png'>" : ""}
            ${blogs[i].typescript ? "<img src='assets/img/ts.png'>" : ""}
          </div>
          <div class="btn-project">
            <button>edit</button>
            <button>delete</button>
          </div>
        </div>`
  } 
}

function getDuration(startDate, endDate) {
  const distance = new Date(endDate) - new Date(startDate)
  const yearDistance = Math.floor(distance / (12 * 30 * 24 * 60 * 60 * 1000))
  if(yearDistance > 0) {
    return yearDistance + " Year"
  } else {
    const monthDistance = Math.floor(distance / (30 * 24 * 60 * 60 * 1000))
    if(monthDistance > 0) {
      return monthDistance + " Month"
    } else {
      const dayDistance = Math.floor(distance / (24 * 60 * 60 * 1000))
      if(dayDistance > 0) {
        return dayDistance + " Day"
      }
    }
  }
}