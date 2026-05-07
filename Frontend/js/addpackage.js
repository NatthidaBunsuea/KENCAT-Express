document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("parcelForm");
  const modal = document.getElementById("successModal");
  const closeBtn = document.querySelector(".close-btn");
  const API_BASE = "http://localhost:8080";

  function getAuthHeaders() {
    const token = localStorage.getItem("token");
    const headers = { "Content-Type": "application/json" };
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    return headers;
  }

  function isValidPhone(phone) {
    return /^\d{9,10}$/.test(phone);
  }

  function isValidWeight(weight) {
    return !isNaN(weight) && Number(weight) > 0;
  }

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const formData = new FormData(form);

    for (const [_, value] of formData.entries()) {
      if (!value.trim()) {
        alert("กรุณากรอกข้อมูลให้ครบถ้วน");
        return;
      }
    }

    if (!isValidPhone(formData.get("senderPhone")) || !isValidPhone(formData.get("receiverPhone"))) {
      alert("เบอร์โทรต้องเป็นตัวเลข 9-10 หลัก");
      return;
    }

    if (!isValidWeight(formData.get("weight"))) {
      alert("น้ำหนักต้องเป็นตัวเลขมากกว่า 0");
      return;
    }

    const payload = {
      deliver: {
        name: formData.get("senderName"),
        surname: formData.get("senderSurname"),
        phone: formData.get("senderPhone"),
        homeNumber: formData.get("senderHomeNumber"),
        soi: formData.get("senderSoi"),
        road: formData.get("senderRoad"),
        district: formData.get("senderDistrict"),
        subdistrict: formData.get("senderSubdistrict"),
        province: formData.get("senderProvince"),
        zipcode: formData.get("senderZipcode")
      },
      receiver: {
        name: formData.get("receiverName"),
        surname: formData.get("receiverSurname"),
        phone: formData.get("receiverPhone"),
        homeNumber: formData.get("receiverHomeNumber"),
        soi: formData.get("receiverSoi"),
        road: formData.get("receiverRoad"),
        district: formData.get("receiverDistrict"),
        subdistrict: formData.get("receiverSubdistrict"),
        province: formData.get("receiverProvince"),
        zipcode: formData.get("receiverZipcode")
      },
      parcel: {
        type: formData.get("delivery"),
        weight: Number(formData.get("weight"))
      }
    };

    try {
      const res = await fetch(`${API_BASE}/api/parcels`, {
        method: "POST",
        headers: getAuthHeaders(),
        credentials: "include",
        body: JSON.stringify(payload)
      });

      const data = await res.json();
      if (!res.ok || data.success === false) {
        alert("❌ เพิ่มพัสดุล้มเหลว: " + (data.message || data.error || "unknown error"));
        return;
      }

      document.getElementById("trackId").innerHTML = `<strong>Track ID :</strong> ${data.trackID}`;
      document.getElementById("parcelId").innerHTML = `<strong>Parcel ID :</strong> ${data.parcelID}`;
      document.getElementById("senderInfo").innerHTML = `<strong>Deliver :</strong> ${formData.get("senderName")} ${formData.get("senderSurname")} ${formData.get("senderPhone")}`;
      document.getElementById("senderAddress").innerHTML = `<strong>Address :</strong> ${formData.get("senderHomeNumber")} ${formData.get("senderSoi")} ${formData.get("senderRoad")} ${formData.get("senderDistrict")} ${formData.get("senderSubdistrict")} ${formData.get("senderProvince")} ${formData.get("senderZipcode")}`;
      document.getElementById("receiverInfo").innerHTML = `<strong>Receiver :</strong> ${formData.get("receiverName")} ${formData.get("receiverSurname")} ${formData.get("receiverPhone")}`;
      document.getElementById("receiverAddress").innerHTML = `<strong>Address :</strong> ${formData.get("receiverHomeNumber")} ${formData.get("receiverSoi")} ${formData.get("receiverRoad")} ${formData.get("receiverDistrict")} ${formData.get("receiverSubdistrict")} ${formData.get("receiverProvince")} ${formData.get("receiverZipcode")}`;

      modal.classList.remove("hidden");
    } catch (err) {
      console.error("❌ Error:", err);
      alert("เกิดข้อผิดพลาดขณะส่งข้อมูลไปเซิร์ฟเวอร์");
    }
  });

  form.addEventListener("reset", () => {
    setTimeout(() => {
      modal.classList.add("hidden");
    }, 100);
  });

  closeBtn.addEventListener("click", () => {
    modal.classList.add("hidden");
  });
});
