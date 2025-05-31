document.addEventListener("DOMContentLoaded", async () => {
    const token = localStorage.getItem("token");

    try {
        let response = await fetch('/admin/dashboard', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        let data = await response.json();

        document.getElementById('totalCampaigns').innerHTML = data.total_campaigns;

        document.getElementById('totalDonations').innerHTML = `$${data.total_donations.toLocaleString()}`;

        document.getElementById('totalUsers').innerHTML = data.total_users;

    } catch (error) {
        console.error('Error fetching dashboard data:', error);
    }
});