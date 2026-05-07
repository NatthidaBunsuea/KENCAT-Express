document.addEventListener("DOMContentLoaded", function () {
  const form = document.querySelector("form");
  const params = new URLSearchParams(window.location.search);
  const trackIdInput = document.getElementById("track-id");
  const presetTrackId = params.get("trackid");
  const API_BASE = "http://localhost:8080";

  if (presetTrackId && trackIdInput) {
    trackIdInput.value = presetTrackId;
  }

  form.addEventListener("submit", async function (e) {
    e.preventDefault();

    const trackID = trackIdInput.value.trim();
    const issue = document.getElementById("issue").value;
    const reason = document.getElementById("reason").value.trim();
    const firstName = document.getElementById("first-name").value.trim();
    const lastName = document.getElementById("last-name").value.trim();
    const phone = document.getElementById("phone").value.trim();
    const email = document.getElementById("email").value.trim();

    const reportData = {
      trackID,
      issue,
      reason,
      firstName,
      lastName,
      phone,
      email
    };

    try {
      const res = await fetch(`${API_BASE}/api/reports`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(reportData)
      });

      const result = await res.json();

      if (!res.ok || result.success === false) {
        alert("❌ ส่งรายงานล้มเหลว: " + (result.message || result.error || "unknown error"));
        return;
      }

      try {
        await emailjs.send("service_3lffs0c", "template_fy17w58", {
          to_email: email,
          to_name: firstName,
          issue,
          reason,
          message: "Kencat Express ได้รับรายงานของคุณแล้ว ขอบคุณสำหรับข้อเสนอแนะของคุณ!"
        });
      } catch (emailErr) {
        console.warn("EmailJS ส่งไม่สำเร็จ แต่บันทึกรายงานเข้า backend แล้ว", emailErr);
      }

      alert("✅ ส่งรายงานสำเร็จ");
      history.back();
    } catch (err) {
      console.error("❌ Error:", err);
      alert("เกิดข้อผิดพลาดในการส่งรายงาน");
    }
  });
});