const API_BASE = "http://localhost:8080";

function buildAuthHeaders() {
  const token = localStorage.getItem("token");
  const headers = { "Content-Type": "application/json" };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
}

window.addEventListener("DOMContentLoaded", async () => {
  const list = document.getElementById("parcel-list");

  try {
    const res = await fetch(`${API_BASE}/api/messenger/tasks`, {
      method: "GET",
      headers: buildAuthHeaders(),
      credentials: "include"
    });
    const parcels = await res.json();

    if (!res.ok) {
      throw new Error(parcels.message || parcels.error || "โหลดงานจัดส่งไม่สำเร็จ");
    }

    parcels.forEach((p) => {
      const item = document.createElement("div");
      item.className = "parcel-item";

      let startDeliveryDisplay = "none";
      let completeBtnDisplay = "none";
      let incompleteBtnDisplay = "none";

      if (p.Status === "Pending") {
        startDeliveryDisplay = "inline-block";
      } else if (p.Status === "In Transit") {
        completeBtnDisplay = "inline-block";
        incompleteBtnDisplay = "inline-block";
      }

      item.innerHTML = `
        <p><strong>Track ID:</strong> ${p.TrackID}</p>
        <p><strong>Parcel ID:</strong> ${p.ParcelID}</p>
        <p><strong>Receiver:</strong> ${p.ReceiverName} ${p.ReceiverSurname}</p>
        <p><strong>Address:</strong> ${p.HomeNumber}, ${p.Soi}, ${p.Road}, ${p.Subdistrict}, ${p.DistrictName}, ${p.ProvinceName} ${p.Zipcode}</p>
        <p><strong>Tel:</strong> ${p.ReceiverTel}</p>
        <p><strong>Status:</strong> <span id="status-${p.ParcelID}">${p.Status}</span></p>
        <p><strong>Deposit Date:</strong> ${new Date(p.DepositDate).toLocaleDateString("th-TH", { year: "numeric", month: "long", day: "numeric" })}</p>
        <button onclick="startDelivery('${p.TrackID}', this, '${p.ParcelID}')" style="display:${startDeliveryDisplay};">ทำการจัดส่ง</button>
        <button id="complete-btn-${p.ParcelID}" style="display:${completeBtnDisplay};" onclick="completeDelivery('${p.TrackID}', '${p.ParcelID}')">จัดส่งสำเร็จ</button>
        <button id="incomplete-btn-${p.ParcelID}" style="display:${incompleteBtnDisplay};" onclick="incompleteDelivery('${p.TrackID}', '${p.ParcelID}')">จัดส่งไม่สำเร็จ</button>
        <hr>
      `;

      list.appendChild(item);
    });
  } catch (err) {
    console.error("❌ Error loading parcels:", err);
    list.innerHTML = `<p>${err.message}</p>`;
  }
});

async function updateTrackingStatus(trackID, status, description, location, parcelID, onSuccess) {
  try {
    const res = await fetch(`${API_BASE}/api/trackings/${encodeURIComponent(trackID)}/status`, {
      method: "PUT",
      headers: buildAuthHeaders(),
      credentials: "include",
      body: JSON.stringify({
        status,
        description,
        location
      })
    });

    const result = await res.json();
    if (!res.ok) {
      throw new Error(result.message || result.error || "อัปเดตสถานะไม่สำเร็จ");
    }

    document.getElementById(`status-${parcelID}`).textContent = status;
    onSuccess();
  } catch (err) {
    console.error("❌ Error updating status:", err);
    alert(err.message);
  }
}

function startDelivery(trackID, button, parcelID) {
  updateTrackingStatus(trackID, "In Transit", "Start delivery", "On route", parcelID, () => {
    alert("เริ่มจัดส่งแล้ว");
    button.style.display = "none";
    document.getElementById(`complete-btn-${parcelID}`).style.display = "inline-block";
    document.getElementById(`incomplete-btn-${parcelID}`).style.display = "inline-block";
  });
}

function completeDelivery(trackID, parcelID) {
  updateTrackingStatus(trackID, "Delivered", "Delivered successfully", "Destination", parcelID, () => {
    alert("จัดส่งสำเร็จ");
    document.getElementById(`complete-btn-${parcelID}`).style.display = "none";
    document.getElementById(`incomplete-btn-${parcelID}`).style.display = "none";
  });
}

function incompleteDelivery(trackID, parcelID) {
  updateTrackingStatus(trackID, "Delivery Failed", "Delivery failed", "Destination", parcelID, () => {
    alert("จัดส่งไม่สำเร็จ");
    document.getElementById(`complete-btn-${parcelID}`).style.display = "none";
    document.getElementById(`incomplete-btn-${parcelID}`).style.display = "none";
  });
}
