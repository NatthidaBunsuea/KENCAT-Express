async function trackParcel() {
  const trackID = document.getElementById("trackingInput").value.trim();
  const resultDiv = document.getElementById("trackResult");

  if (!trackID) {
    resultDiv.innerHTML = `<span style="color:red;">❌ กรุณากรอก Track ID ก่อน</span>`;
    return;
  }

  resultDiv.innerHTML = "🔍 กำลังค้นหา...";

  try {
    const res = await fetch(`http://localhost:8080/api/trackings/${encodeURIComponent(trackID)}`);
    if (!res.ok) throw new Error("ไม่พบพัสดุ");

    const data = await res.json();
    const status = (data.Status || "").trim();

    if (status === "Delivered") {
      window.location.href = `html/package_status_success.html?trackid=${encodeURIComponent(trackID)}`;
      return;
    }

    if (status === "Delivery Failed") {
      window.location.href = `html/package_status_failed.html?trackid=${encodeURIComponent(trackID)}`;
      return;
    }

    if (status === "Pending" || status === "In Transit") {
      window.location.href = `html/package_status_inprocess.html?trackid=${encodeURIComponent(trackID)}`;
      return;
    }

    resultDiv.innerHTML = `<span style="color:red;">❌ ไม่รองรับสถานะ: ${status || "ไม่มีข้อมูล"}</span>`;
  } catch (err) {
    resultDiv.innerHTML = `<span style="color:red;">❌ ${err.message}</span>`;
  }
}
