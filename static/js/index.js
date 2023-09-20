const arrayContainerProdutos = document.querySelectorAll('.outer-container')
const arrayProdutos = document.querySelectorAll('.inner-container')

arrayContainerProdutos.forEach((element, index) => {
   element.addEventListener('click', ()=>{
    arrayProdutos[index].classList.toggle('hidden')
   }) 
});