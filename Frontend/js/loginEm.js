function showLoginFailModal(errorMessage) {
    document.getElementById('errorText').textContent = errorMessage;
    const modal = document.getElementById('loginFailModal');
    modal.classList.add('show');
}

function closeModal() {
    const modal = document.getElementById('loginFailModal');
    modal.classList.remove('show');
}

document.addEventListener('DOMContentLoaded', () => {
    const closeBtn = document.querySelector('.close-btn');
    if (closeBtn) {
        closeBtn.addEventListener('click', closeModal);
    }

    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const employeeID = document.getElementById('employeeID').value.trim();
        const password = document.getElementById('password').value.trim();
        const role = document.querySelector('input[name="role"]:checked')?.value || '';

        document.getElementById('employeeID-error').textContent = '';
        document.getElementById('password-error').textContent = '';

        let isValid = true;

        if (!employeeID) {
            document.getElementById('employeeID-error').textContent = 'Please enter Employee ID.';
            isValid = false;
        }

        if (!password) {
            document.getElementById('password-error').textContent = 'Please enter Password.';
            isValid = false;
        }

        if (!role) {
            await Swal.fire({
                icon: 'warning',
                title: 'Please select a role!',
                confirmButtonColor: '#3085d6'
            });
            isValid = false;
        }

        if (!isValid) return;

        try {
            const res = await fetch('http://localhost:8080/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ employeeID, password, role })
            });

            const data = await res.json();

            if (!res.ok || data.success === false) {
                showLoginFailModal(data.message || data.error || 'Invalid credentials.');
                return;
            }

            const token = data.token || '';
            const employee = data.user || null;

            if (token) localStorage.setItem('token', token);
            if (employee) localStorage.setItem('employee', JSON.stringify(employee));
            localStorage.setItem('employeeRoleSelected', role);

            await Swal.fire({
                icon: 'success',
                title: 'Login Successful!',
                text: `Welcome ${role}!`,
                confirmButtonColor: '#3085d6'
            });

            if (role === 'Parcel Clerk') {
                window.location.href = 'packageList_ParcelClerk.html';
            } else {
                window.location.href = 'packageSend_Messenger.html';
            }
        } catch (err) {
            console.error(err);
            showLoginFailModal('Cannot connect to server.');
        }
    });
});