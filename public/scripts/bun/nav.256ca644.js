// @bun
console.log("nav js loaded");var o=document.querySelectorAll("#nav a");o.forEach((e)=>{console.log("adding event listener for link: ",e),e.addEventListener("click",(l)=>{console.log("click"),l.preventDefault(),console.log(`Clicked link: ${e}`)})});
