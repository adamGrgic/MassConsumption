

console.log('nav js loaded');

const links = document.querySelectorAll("#navbar a")
console.log(links);
links.forEach(link => {
    console.log('adding event listener for link: ',link)
    link.addEventListener('click', (event) => {
        console.log('click');
        event.preventDefault();
        console.log(`Clicked link: ${link}`);
        console.log(event);
        links.forEach(l => l.removeAttribute('active'));
        link.setAttribute('active', 'true');
    });
  });
