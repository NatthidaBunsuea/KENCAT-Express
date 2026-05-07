console.log("Loaded packageinfo.js");

function formatThaiDate(value) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("th-TH");
}

function buildTimeline(data) {
  const eventItems = Array.isArray(data.Events)
    ? data.Events
        .slice()
        .sort((a, b) => new Date(b.createdAt || b.CreatedAt || 0) - new Date(a.createdAt || a.CreatedAt || 0))
        .map((event) => {
          const createdAt = event.createdAt || event.CreatedAt;
          const status = event.status || event.Status || "-";
          const description = event.description || event.Description || "";
          const location = event.location || event.Location || "";
          return `<div class="timeline-item"><strong>${formatThaiDate(createdAt)}</strong> - ${status}${description ? ` : ${description}` : ""}${location ? ` (${location})` : ""}</div>`;
        })
        .join("")
    : "";

  return `
    <p>สถานะล่าสุด: ${data.Status || "ไม่มีข้อมูล"}</p>
    <p>ชื่อผู้ส่ง: ${data.Sender || "ไม่มีข้อมูลผู้ส่ง"}</p>
    <p>ชื่อผู้รับ: ${data.Receiver || "ไม่มีข้อมูลผู้รับ"}</p>
    <p>ที่อยู่ผู้รับ: ${data.Address || "ไม่มีข้อมูลที่อยู่"}</p>
    <p>ชื่อพนักงานจัดส่ง: ${data.Deliverer || "ยังไม่มีข้อมูล"}</p>
    <p>ประเภทรถ: ${data.Typecar || "ยังไม่มีข้อมูล"}</p>
    <p>เลขทะเบียนรถ: ${data.License || "ยังไม่มีข้อมูล"}</p>
    <p>วันที่ส่งถึง: ${formatThaiDate(data.DeliveredAt)}</p>
    <p>วันที่อัปเดตล่าสุด: ${formatThaiDate(data.UpdatedAt)}</p>
    ${eventItems || ""}
  `;
}

window.addEventListener("DOMContentLoaded", async () => {
  const params = new URLSearchParams(window.location.search);
  const trackID = params.get("trackid");
  const timeline = document.querySelector(".timeline");
  const reportButton = document.querySelector(".btn");

  if (!trackID) return;

  if (reportButton) {
    reportButton.onclick = () => {
      window.location.href = `package_report.html?trackid=${encodeURIComponent(trackID)}`;
    };
  }

  try {
    const res = await fetch(`http://localhost:8080/api/trackings/${encodeURIComponent(trackID)}`);
    if (!res.ok) throw new Error("ไม่พบข้อมูล");

    const data = await res.json();
    document.querySelector(".status-header h3").textContent = `เลขพัสดุ (Track ID): ${trackID}`;

    const deposit = formatThaiDate(data.UpdatedAt);
    const delivered = data.DeliveredAt ? formatThaiDate(data.DeliveredAt) : data.Status || "-";
    const headerParagraph = document.querySelectorAll(".status-header p")[0];
    if (headerParagraph) {
      headerParagraph.textContent = `ระหว่างวันที่ ${deposit} - ${delivered}`;
    }

    document.body.classList.remove("status-success", "status-failed", "status-pending");

    const steps = document.querySelectorAll(".step");
    if (steps.length >= 3) {
      if (data.Status === "Delivered") {
        document.body.classList.add("status-success");
        steps[1].innerHTML = "2<br>อยู่ระหว่างจัดส่ง";
        steps[2].innerHTML = "3<br>จัดส่งสำเร็จแล้ว";
      } else if (data.Status === "Delivery Failed") {
        document.body.classList.add("status-failed");
        steps[1].innerHTML = "2<br>อยู่ระหว่างจัดส่ง";
        steps[2].innerHTML = "3<br>จัดส่งไม่สำเร็จ";
      } else if (data.Status === "In Transit") {
        document.body.classList.add("status-pending");
        steps[1].innerHTML = "2<br>อยู่ระหว่างจัดส่ง";
        steps[2].innerHTML = "3<br>รอผลการจัดส่ง";
      } else {
        document.body.classList.add("status-pending");
        steps[1].innerHTML = "2<br>รอจัดส่ง";
        steps[2].innerHTML = "3<br>ยังไม่ถึงขั้นตอนนี้";
      }
    }

    timeline.innerHTML = buildTimeline(data);
  } catch (err) {
    console.error("เกิดข้อผิดพลาด:", err);
    timeline.innerHTML = `<p>ไม่พบข้อมูลพัสดุ</p>`;
  }
});
