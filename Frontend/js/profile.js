let userRole = "";

function readEmployeeFromStorage() {
  try {
    return JSON.parse(localStorage.getItem("employee") || "null");
  } catch (err) {
    console.error("Cannot parse employee from localStorage", err);
    return null;
  }
}

function resolveRoleCode(employee) {
  const explicitRole = (employee?.role || localStorage.getItem("employeeRoleSelected") || "").toString().trim().toUpperCase();
  if (explicitRole === "PARCEL CLERK" || explicitRole === "PARCEL_CLERK" || explicitRole === "P") return "P";
  if (explicitRole === "MESSENGER" || explicitRole === "M") return "M";
  if (explicitRole === "ADMIN" || explicitRole === "A") return "A";
  return "";
}

function loadProfile() {
  const employee = readEmployeeFromStorage();
  if (!employee) {
    alert("กรุณาเข้าสู่ระบบก่อน");
    window.location.href = "/html/loginEm.html";
    return;
  }

  userRole = resolveRoleCode(employee);

  document.getElementById("EmployeeID").textContent = employee.employeeCode || employee.id || "-";
  document.getElementById("role").textContent = employee.displayRole || employee.role || localStorage.getItem("employeeRoleSelected") || "-";
  document.getElementById("name").textContent = `${employee.firstName || ""} ${employee.lastName || ""}`.trim() || "-";
  document.getElementById("email").textContent = employee.email || "-";
  document.getElementById("birthday").textContent = "-";

  const roleButton = document.querySelector(".role-btn");
  if (roleButton) {
    roleButton.disabled = false;
  }
}

function redirectToRole() {
  if (userRole === "P") {
    window.location.href = "/html/packageList_ParcelClerk.html";
  } else if (userRole === "M") {
    window.location.href = "/html/packageSend_Messenger.html";
  } else if (userRole === "A") {
    window.location.href = "/html/packageList_ParcelClerk.html";
  } else {
    alert("ไม่พบบทบาทผู้ใช้ที่ถูกต้อง");
    console.log("userRole คือ:", userRole);
  }
}

function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("employee");
  localStorage.removeItem("employeeRoleSelected");
  alert("ออกจากระบบแล้ว (ล้าง session ฝั่งหน้าเว็บ)");
  window.location.href = "/html/loginEm.html";
}

document.addEventListener("DOMContentLoaded", loadProfile);
