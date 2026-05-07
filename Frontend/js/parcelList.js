const API_BASE = "http://localhost:8080";

function buildAuthHeaders() {
  const token = localStorage.getItem("token");
  const headers = {};
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
}

window.addEventListener("DOMContentLoaded", async () => {
  await loadParcels();
});

async function loadParcels() {
  const list = document.getElementById("parcel-list");
  list.innerHTML = "";

  try {
    const res = await fetch(`${API_BASE}/api/parcels`, {
      method: "GET",
      headers: buildAuthHeaders(),
      credentials: "include"
    });

    const parcels = await res.json();
    if (!res.ok) {
      throw new Error(parcels.message || parcels.error || "โหลดข้อมูลไม่สำเร็จ");
    }

    window.parcelsData = parcels;
    renderParcels(parcels);
  } catch (err) {
    console.error("❌ Error loading parcels:", err);
    list.innerHTML = `<p>${err.message}</p>`;
  }
}

function renderParcels(parcels) {
  const list = document.getElementById("parcel-list");
  list.innerHTML = "";

  parcels.forEach((p) => {
    const item = document.createElement("div");
    item.className = "parcel-item";

    const fullName = `${p.ReceiverName} ${p.ReceiverSurname}`;
    const fullAddress = `${p.HomeNumber}, ${p.Soi || ""}, ${p.Road || ""}, ${p.Subdistrict}, ${p.DistrictName}, ${p.ProvinceName} ${p.Zipcode}`;

    item.innerHTML = `
      <div><strong>Track ID:</strong> <span>${p.TrackID}</span></div>
      <div><strong>Parcel ID:</strong> <span>${p.ParcelID}</span></div>
      <div><strong>Receiver:</strong> <span>${fullName}</span></div>
      <div><strong>Address:</strong> <span>${fullAddress}</span></div>
      <div><strong>Tel:</strong> <span>${p.ReceiverTel}</span></div>
      <div><strong>Status:</strong> <span>${p.Status}</span></div>
      <div><strong>Deposit Date:</strong> <span>${formatDate(p.DepositDate)}</span></div>
      <hr>
    `;

    list.appendChild(item);
  });
}

function searchParcel() {
  const keyword = (document.getElementById("searchInput")?.value || "").trim().toLowerCase();
  const data = Array.isArray(window.parcelsData) ? window.parcelsData : [];
  if (!keyword) {
    renderParcels(data);
    return;
  }

  const filtered = data.filter((parcel) => {
    return [parcel.TrackID, parcel.ParcelID, parcel.ReceiverName, parcel.ReceiverSurname, parcel.Status]
      .filter(Boolean)
      .some((value) => value.toString().toLowerCase().includes(keyword));
  });

  renderParcels(filtered);
}

function formatDate(dateString) {
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("th-TH");
}
