function predictArea() {
    let year = document.getElementById("year").value;
    fetch(`/predict?year=${year}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("result").innerText =
                `Predicted urban area for ${data.year}: ${data.predicted_area.toFixed(2)} sq km`;
        });
}
