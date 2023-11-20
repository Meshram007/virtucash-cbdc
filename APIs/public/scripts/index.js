/*
const addData = async (event) => {
    event.preventDefault();

    const carId = document.getElementById("carId").value;
    const make = document.getElementById("make").value;
    const model = document.getElementById("model").value;
    const color = document.getElementById("color").value;
    const dateOfManufacture = document.getElementById("dateOfManufacture").value;
    const manufacturerName = document.getElementById("manufacturerName").value;


    const carData = {
        carId: carId,
        make: make,
        model: model,
        color: color,
        dateOfManufacture: dateOfManufacture,
        manufacturerName: manufacturerName,
    };


    if (
        carId.length == 0 ||
        make.length == 0 ||
        model.length == 0 ||
        color.length == 0 ||
        dateOfManufacture.length == 0 ||
        manufacturerName.length == 0
    ) {
        alert("Please enter the data properly.");
    } else {
        try {
            const response = await fetch("/api/car", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(carData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            return alert("Car Created");
        } catch (err) {
            alert("Error");
            console.log(err);
        }
    }
};

const readData = async (event) => {
    event.preventDefault();
    const carId = document.getElementById("carIdInput").value;

    if (carId.length == 0) {
        alert("Please enter a valid ID.");
    } else {
        try {
            const response = await fetch(`/api/car/${carId}`);
            let responseData = await response.json();
            console.log("sadsdsdd", responseData);
            alert(JSON.stringify(responseData));
        } catch (err) {
            alert("Error");
            console.log(err);
        }
    }
};
*/
////////

// function toManuDash() {
//     window.location.href = '/manufacturer';
// }

function swalBasic(data) {
    swal.fire({
        // toast: true,
        icon: `${data.icon}`,
        title: `${data.title}`,
        animation: true,
        position: 'center',
        showConfirmButton: true,
        footer: `${data.footer}`,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', swal.stopTimer)
            toast.addEventListener('mouseleave', swal.resumeTimer)
        }
    });
}

function swalBasicRefresh(data) {
    swal.fire({
        // toast: true,
        icon: `${data.icon}`,
        title: `${data.title}`,
        animation: true,
        position: 'center',
        showConfirmButton: true,
        footer: `${data.footer}`,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', swal.stopTimer)
            toast.addEventListener('mouseleave', swal.resumeTimer)
        }
    }).then(() => {
        location.reload();
    });
}

function reloadWindow() {
    window.location.reload();
}

const addCarData = async (event) => {
    event.preventDefault();
    const carId = document.getElementById('carId').value;
    const make = document.getElementById('carMake').value;
    const model = document.getElementById('carModel').value;
    const color = document.getElementById('carColour').value;
    const dateOfManufacture = document.getElementById('dom').value;
    const manufacturerName = document.getElementById('manufacturer').value;
    console.log(carId + make + model + color);

    const carData = {
        carId: carId,
        make: make,
        model: model,
        color: color,
        dateOfManufacture: dateOfManufacture,
        manufacturerName: manufacturerName,
    };
    if (
        carId.length == 0 ||
        make.length == 0 ||
        model.length == 0 ||
        color.length == 0 ||
        dateOfManufacture.length == 0 ||
        manufacturerName.length == 0
    ) {
        const data = {
            title: "You might have missed something",
            footer: "Enter all mandatory fields to add a new car",
            icon: "warning"
        }
        swalBasic(data);
    } else {
        try {
            const response = await fetch("/api/car", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(carData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const carStatus = {
                title: "Success",
                footer: "Added a new car",
                icon: "success"
            }
            swalBasicRefresh(carStatus);

        } catch (err) {
            // alert("Error");
            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }

}

const readCarData = async (event) => {
    event.preventDefault();
    const carId = document.getElementById("queryCarId").value;
    if (carId.length == 0) {
        const data = {
            title: "Enter a valid car Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`/api/car/${carId}`);
            let responseData = await response.json();
            console.log("response", responseData);
            // alert(JSON.stringify(responseData));
            const dataBuf = JSON.stringify(responseData)
            swal.fire({
                // toast: true,
                icon: `success`,
                title: `Current status of car with carId ${carId} :`,
                animation: false,
                position: 'center',
                html: `<h3>${dataBuf}</h3>`,
                showConfirmButton: true,
                timer: 3000,
                timerProgressBar: true,
                didOpen: (toast) => {
                    toast.addEventListener('mouseenter', swal.stopTimer)
                    toast.addEventListener('mouseleave', swal.resumeTimer)
                }
            })
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};

function getMatchingOrders(carId) {
    // console.log("carId", carId)
    window.location.href = '/api/order/match-car?carId=' + carId;
}

//Method to get the history of an item
function getCarHistory(carId) {
    console.log("carId====", carId)
    window.location.href = '/api/car/history?carId=' + carId;
}


const registerCar = async (event) => {
    // function registerCar(event) {
    console.log("Entered the register function")
    event.preventDefault();
    const carId = document.getElementById('carId').value;
    const carOwner = document.getElementById('carOwner').value;
    const regNumber = document.getElementById('regNumber').value;
    console.log(carId + carOwner + regNumber);
    const carData = {
        carId: carId,
        carOwner: carOwner,
        regNumber: regNumber,
    };
    if (carId.length == 0 || carOwner.length == 0 || regNumber.length == 0) {
        const data = {
            title: "You have missed something",
            footer: "All fields are mandatory",
            icon: "warning"
        }
        swalBasic(data)
    }
    else {
        try {
            const response = await fetch("/api/car/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(carData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const carStatus = {
                title: `Registered car ${carId} to ${carOwner}`,
                footer: "Registered car",
                icon: "success"
            }
            swalBasicRefresh(carStatus);

        } catch (err) {
            console.log(err);
            const data = {
                title: `Failed to register car`,
                footer: "Please try again !!",
                icon: "error"
            }
            swalBasic(data);
        }


    }
}


const addOrder = async (event) => {
    event.preventDefault();
    const orderNumber = document.getElementById('orderNumber').value;
    const carMake = document.getElementById('carMake').value;
    const carModel = document.getElementById('carModel').value;
    const carColour = document.getElementById('carColour').value;
    const dealerName = document.getElementById('dealerName').value;
    console.log(orderNumber + carColour + dealerName);

    const orderData = {
        orderId: orderNumber,
        make: carMake,
        model: carModel,
        color: carColour,
        dealerName: dealerName,
    };
    if (
        orderNumber.length == 0 ||
        carMake.length == 0 ||
        carModel.length == 0 ||
        carColour.length == 0 ||
        dealerName.length == 0
    ) {
        const data = {
            title: "You have missed something",
            footer: "All fields are mandatory",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch("/api/order", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(orderData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            // return alert("Order Created");
            const orderStatus = {
                title: `Order is created`,
                footer: "Raised Order",
                icon: "success"
            }
            swalBasic(orderStatus)
        } catch (err) {
            // alert("Error");
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
            console.log(err);
        }
    }

}

const readOrder = async (event) => {
    event.preventDefault();
    const orderId = document.getElementById("ordNum").value;

    if (orderId.length == 0) {
        const data = {
            title: "Enter a order number",
            footer: "Order Number is mandatory",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`/api/order/${orderId}`);
            let responseData = await response.json();
            console.log("response", responseData);
            const dataBuf = JSON.stringify(responseData)
            swal.fire({
                // toast: true,
                icon: `success`,
                title: `Current status of Order : `,
                animation: false,
                position: 'center',
                html: `<h3>${dataBuf}</h3>`,
                showConfirmButton: true,
                timer: 3000,
                timerProgressBar: true,
                didOpen: (toast) => {
                    toast.addEventListener('mouseenter', swal.stopTimer)
                    toast.addEventListener('mouseleave', swal.resumeTimer)
                }
            })
            // alert(JSON.stringify(responseData));
        } catch (err) {
            // alert("Error");
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
            console.log(err);
        }
    }
};


async function matchOrder(orderId, carId) {
    if (!orderId || !carId) {
        const data = {
            title: "Enter a order number",
            footer: "Order Number is mandatory",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        const matchData = {
            carId: carId,
            orderId: orderId,
        }
        try {
            const response = await fetch("/api/car/match-order", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(matchData),
            });
            const data = {
                title: `Order matched successfully`,
                footer: "Order matched",
                icon: "success"
            }
            swalBasicRefresh(data)

        } catch (err) {
            const data = {
                title: `Failed to match order`,
                footer: "Please try again !!",
                icon: "error"
            }
            swalBasic(data)
        }


    }
}


function allOrders() {
    window.location.href = '/api/order/all';
}


async function getEvent() {
    try {
        const response = await fetch("/api/event", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });
        // console.log("RESPONSE: ", response)
        const data = await response.json()
        // console.log("DATA: ", data);

        const eventsData = data["carEvent"]
        swal.fire({
            toast: true,
            // icon: `${data.icon}`,
            title: `Event : `,
            animation: false,
            position: 'top-right',
            html: `<h5>${eventsData}</h5>`,
            showConfirmButton: false,
            timer: 5000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', swal.stopTimer)
                toast.addEventListener('mouseleave', swal.resumeTimer)
            }
        })
    } catch (err) {
        swal.fire({
            toast: true,
            icon: `error`,
            title: `Error`,
            animation: false,
            position: 'top-right',
            showConfirmButton: true,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', swal.stopTimer)
                toast.addEventListener('mouseleave', swal.resumeTimer)
            }
        })
        console.log(err);
    }
}





