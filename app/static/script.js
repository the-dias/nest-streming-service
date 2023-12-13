'use strict';

const button = document.getElementById("btn_id_input")
const orders = document.getElementsByClassName("display_order")
const order_id = document.getElementById('display_order_id')
const data = document.getElementById("jsonData")



if(data != null) {

    const formattedJson = JSON.stringify(JSON.parse(data.textContent), null, 2);

    document.getElementById("showDataJson").innerHTML = formattedJson

    console.log(formattedJson)
}


if(order_id.textContent.indexOf('-1') != -1 ) {
    for (let item of orders) {
        item.style.display = 'none';
    }
}


button.addEventListener('click', function(e) {
    const input = document.getElementById("id_input")
    console.log('input: ' ,input.value)
    const inputValue = input.value
    e.preventDefault()
    console.log('input value: ', inputValue)
    if(inputValue === '') {
        window.location.href = "http://localhost:3000/order";
        return;
    }
    const URL = `http://localhost:3000/order?id=${+inputValue}`

    console.log(order_id.textContent)

    window.location.href = URL;
    
})