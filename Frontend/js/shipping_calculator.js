document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form");
  const weightInput = document.getElementById("weight");
  const originSelect = document.getElementById("origin");
  const destinationSelect = document.getElementById("destination");
  const API_BASE = "http://localhost:8080";

  const locations = [
    "Bangkok",
    "Nonthaburi",
    "Pathum Thani",
    "Samut Prakan",
    "Chiang Mai",
    "Phuket",
    "Khon Kaen",
    "Songkhla"
  ];

  for (const city of locations) {
    const option1 = document.createElement("option");
    option1.value = city;
    option1.textContent = city;
    originSelect.appendChild(option1);

    const option2 = document.createElement("option");
    option2.value = city;
    option2.textContent = city;
    destinationSelect.appendChild(option2);
  }

  form.addEventListener("submit", async function (e) {
    e.preventDefault();

    const origin = originSelect.value;
    const destination = destinationSelect.value;
    const weight = parseFloat(weightInput.value);
    const shippingType = document.querySelector("input[name='shipping_type']:checked");
    const resultDiv = document.getElementById("result");

    if (!origin || !destination || Number.isNaN(weight) || weight <= 0 || !shippingType) {
      Swal.fire("กรุณากรอกข้อมูลให้ครบทุกช่องและน้ำหนักต้องมากกว่า 0");
      return;
    }

    const deliveryType = shippingType.id === "express" ? "express" : "standard";
    const query = new URLSearchParams({
      origin,
      destination,
      deliveryType,
      weight: weight.toString()
    });

    try {
      const res = await fetch(`${API_BASE}/api/shipping/calculate?${query.toString()}`);
      const result = await res.json();
      if (!res.ok) {
        throw new Error(result.message || result.error || "คำนวณค่าส่งไม่สำเร็จ");
      }

      const totalCost = Number(result.shippingCost || 0).toFixed(2);
      Swal.fire({
        title: '<span style="font-size: 0.9em;">คำนวณค่าส่งแล้ว!</span>',
        html: `<div style="font-size: 1.4em; font-weight: bold;">ค่าจัดส่งคือ ${totalCost} บาท</div>`,
        icon: "success",
        confirmButtonText: "ตกลง"
      }).then(() => {
        resultDiv.textContent = `ค่าจัดส่งโดยประมาณ: ${totalCost} บาท`;
      });
    } catch (err) {
      console.error(err);
      Swal.fire("เกิดข้อผิดพลาด", err.message, "error");
    }
  });

  weightInput.addEventListener("input", () => {
    weightInput.value = weightInput.value.replace(/[^\d.]/g, "");
  });
});
