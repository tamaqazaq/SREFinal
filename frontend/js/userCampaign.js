async function fetchUserCampaigns() {
    const token = localStorage.getItem("token");
    const decoded = jwt_decode(token);
    const userId = decoded.UserID;

    if (!userId) {
        console.error("User ID not found");
        return;
    }

    const response = await fetch(`/campaigns/user/${userId}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    const data = await response.json();
    const campaignList = document.getElementById("user-campaign-list");
    campaignList.innerHTML = "";

    if (!data || data.length === 0) {
        campaignList.innerHTML = "<p>You have no campaigns.</p>";
        return;
    }

    data.forEach(campaign => {
        const card = document.createElement("div");
        card.classList.add("campaign-card", "clickable-card");
        card.innerHTML = `
  <div class="card-content">
    <h3>${campaign.title}</h3>
    <div class="campaign-media">
      <img src="${campaign.media_path}" alt="Campaign Image">
    </div>
    <p>${campaign.description}</p>
    <p><strong>Campaign ID:</strong> ${campaign.campaign_id}</p>
    <p><strong>Target:</strong> $${campaign.target_amount}</p>
    <p><strong>Raised:</strong> $${campaign.amount_raised}</p>
    <p><strong>Status:</strong> ${campaign.status}</p>
    <div class="social-buttons">
      <!-- Twitter Share Button -->
      <a href="https://twitter.com/intent/tweet?url=https://yourcrowdfundingurl.com/campaigns/${campaign.campaign_id}"
         class="btn share-twitter" target="_blank">
         <svg viewBox="0 0 24 24" width="20" height="20">
           <path d="M23.954 4.569c-.885.389-1.83.654-2.825.775
                    1.014-.611 1.794-1.574 2.163-2.723-.951.555-2.005.959-3.127 1.184
                    -.896-.957-2.173-1.555-3.591-1.555-3.42 0-5.966 3.156-5.175 6.533
                    -4.091-.205-7.719-2.165-10.148-5.144-1.29 2.214-.669 5.108 1.523 6.574
                    -.807-.026-1.566-.248-2.228-.616-.054 2.281 1.581 4.415 3.949 4.89
                    -.693.188-1.452.232-2.224.084.63 1.953 2.445 3.376 4.604 3.418
                    -2.07 1.625-4.678 2.348-7.29 2.04 2.179 1.394 4.768 2.214 7.557 2.214
                    9.142 0 14.307-7.721 13.995-14.646.961-.695 1.796-1.562 2.457-2.549z"/>
         </svg>
         Twitter
      </a>
      <!-- LinkedIn Share Button -->
      <a href="https://www.linkedin.com/shareArticle?mini=true&url=https://yourcrowdfundingurl.com/campaigns/${campaign.campaign_id}"
         class="btn share-linkedin" target="_blank">
         <svg viewBox="0 0 24 24" width="20" height="20">
           <path d="M22.23 0H1.77C.79 0 0 .775 0 1.732V22.27c0
                    .956.79 1.73 1.77 1.73h20.46c.98 0 1.77-.775
                    1.77-1.73V1.732C24 .775 23.21 0 22.23 0zM7.05
                    20.452H3.56V9h3.49v11.452zM5.3 7.682c-1.11 0-2.01-.91-2.01-2.033
                    0-1.12.9-2.034 2.01-2.034s2.01.914 2.01 2.034c0 1.122-.9 2.033-2.01
                    2.033zm14.15 12.77h-3.48V14.34c0-1.453-.03-3.325-2.025-3.325-2.03 0-2.34
                    1.584-2.34 3.22v6.217h-3.48V9h3.34v1.57h.05c.47-.9
                    1.61-1.848 3.31-1.848 3.55 0 4.21 2.34 4.21
                    5.383v6.347z"/>
         </svg>
         LinkedIn
      </a>
    </div>
    <div class="campaign-actions">
      <button class="btn-edit" onclick="openEditModal(
        ${campaign.campaign_id},
        '${escapeQuotes(campaign.title)}',
        '${escapeQuotes(campaign.description)}',
        ${campaign.target_amount},
        '${campaign.status}', '${campaign.category}'
      ); event.stopPropagation();">Edit</button>
      <button class="btn-delete" onclick="deleteCampaign(${campaign.campaign_id}); event.stopPropagation();">Delete</button>
    </div>
  </div>
`;

// Make the entire card clickable for details (except when inner buttons are clicked)
        card.addEventListener("click", () => {
            window.location.href = `/static/campaign-details.html?id=${campaign.campaign_id}`;
        });

        campaignList.appendChild(card);
    });

    document.querySelectorAll(".share-facebook").forEach(button => {
        button.addEventListener("click", async (e) => {
            e.preventDefault();
            const url = button.dataset.url;
            const response = await fetch(`/share?url=${encodeURIComponent(url)}&text=Check out this campaign!`);
            const data = await response.json();
            window.open(data.facebook, '_blank');
        });
    });

    document.querySelectorAll(".share-twitter").forEach(button => {
        button.addEventListener("click", async (e) => {
            e.preventDefault();
            const url = button.dataset.url;
            const response = await fetch(`/share?url=${encodeURIComponent(url)}&text=Check out this campaign!`);
            const data = await response.json();
            window.open(data.twitter, '_blank');
        });
    });

    document.querySelectorAll(".share-linkedin").forEach(button => {
        button.addEventListener("click", async (e) => {
            e.preventDefault();
            const url = button.dataset.url;
            const response = await fetch(`/share?url=${encodeURIComponent(url)}&text=Check out this campaign!`);
            const data = await response.json();
            window.open(data.linkedin, '_blank');
        });
    });
}

async function createCampaign() {
    const title = document.getElementById('create-title').value;
    const description = document.getElementById('create-description').value;
    const targetAmount = document.getElementById('create-target-amount').value;
    const status = document.getElementById('create-status').value;
    const mediaFile = document.getElementById('create-media').files[0];
    const token = localStorage.getItem("token");
    const decoded = jwt_decode(token);
    const userId = decoded.UserID;
    const category = document.getElementById('create-category').value;

    if (!title || !description || !targetAmount || !status || !category) {
        alert("Please fill in all fields.");
        return;
    }
    if (!userId) {
        console.error("User ID not found");
        return;
    }

    let mediaPath = null;

    if (mediaFile) {

        const formData = new FormData();
        formData.append("media", mediaFile);

        try {
            const uploadResponse = await fetch('/upload', {
                method: 'POST',
                body: formData,
            });

            const uploadData = await uploadResponse.json();
            if (!uploadData.success) {
                alert("Failed to upload file");
                return;
            }

            mediaPath = uploadData.file_url; // ✅ Store uploaded file URL
        } catch (error) {
            console.error("Error uploading file:", error);
            alert("File upload failed.");
            return;
        }
    }


    const campaignData = {
        user_id: userId,
        title: title,
        description: description,
        target_amount: parseFloat(targetAmount),
        status: status,
        category: category,
        media_path: mediaPath,
    };

    try {
        const response = await fetch('/campaigns/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(campaignData),
        });

        const data = await response.json();
        if (data.success) {
            alert('Campaign created successfully!');
            hideCreateForm();
            fetchUserCampaigns();
        } else {
            alert(`Failed to create campaign: ${data.message}`);
        }
    } catch (error) {
        console.error('Error creating campaign:', error);
    }
}
function showCreateForm() {
    document.getElementById('create-form').style.display = 'block';
}
function escapeQuotes(str) {
    if (!str) return "";
    return str.replace(/'/g, "\\'").replace(/"/g, '\\"');
}
function hideCreateForm() {
    document.getElementById('create-form').style.display = 'none';
}
document.getElementById("editCampaignForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    const campaignId = document.getElementById("editCampaignId").value;
    const newTitle = document.getElementById("editTitle").value;
    const newDescription = document.getElementById("editDescription").value;
    const newTargetAmount = document.getElementById("editTargetAmount").value;
    const status = document.getElementById("editStatus").value;
    const mediaFile = document.getElementById("editMedia").files[0];
    const token = localStorage.getItem("token");
    const category = document.getElementById('editCategory').value;
    const targetAmountNum = parseFloat(newTargetAmount);

    // Create a FormData object (which supports file uploads)
   /* const formData = new FormData();
    formData.append("title", newTitle);
    formData.append("description", newDescription);
    formData.append("target_amount", targetAmountNum);
    formData.append("status", status);*/
    let mediaPath = null;
    if (mediaFile) {

        const formData = new FormData();
        formData.append("media", mediaFile);

        try {
            const uploadResponse = await fetch('/upload', {
                method: 'POST',
                body: formData,
            });

            const uploadData = await uploadResponse.json();
            if (!uploadData.success) {
                alert("Failed to upload file");
                return;
            }

            mediaPath = uploadData.file_url;
        } catch (error) {
            console.error("Error uploading file:", error);
            alert("File upload failed.");
            return;
        }
    }
   else if(mediaFile == "") {
       mediaPath = "";
    }

    const campaignData = {
        title: newTitle,
        description: newDescription,
        target_amount: parseFloat(newTargetAmount),
        status: status,
        category: category,
        media_path: mediaPath,
    };


    try {
        const response = await fetch(`/campaigns/${campaignId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(campaignData),
        });

        const result = await response.json();
        if (result.success) {
            alert("Campaign updated successfully!");
            // Hide the modal
            const modalEl = document.getElementById('editCampaignModal');
            const modalInstance = bootstrap.Modal.getInstance(modalEl);
            modalInstance.hide();
            fetchUserCampaigns();  // Refresh the campaign list
        } else {
            alert(`Failed to update campaign: ${result.message || "Unknown error"}`);
        }
    } catch (error) {
        console.error("Error updating campaign:", error);
        alert("An error occurred while updating the campaign.");
    }
});

/**
 * Opens the edit modal with the campaign data pre-filled.
 * @param {number} id - The campaign ID.
 * @param {string} title - The current title.
 * @param {string} description - The current description.
 * @param {number} targetAmount - The current target amount.
 * @param {string} status - The current status.
 * @param {string} category - The current category.
 */
function openEditModal(id, title, description, targetAmount, status, category ) {
    document.getElementById('editCampaignId').value = id;
    document.getElementById('editTitle').value = title;
    document.getElementById('editDescription').value = description;
    document.getElementById('editTargetAmount').value = targetAmount;
    document.getElementById('editStatus').value = status;
  //  document.getElementById('editMedia').value = ""; // Clear any previous file selection
   document.getElementById('editCategory').value = category
    // Use Bootstrap's modal method to show the modal
    const editModal = new bootstrap.Modal(document.getElementById('editCampaignModal'));
    editModal.show();
}
fetchUserCampaigns();

async function deleteCampaign(campaign_id) {
        if (!confirm("Are you sure you want to delete this campaign?")) return;

        const response = await fetch(`/campaigns/${campaign_id}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });

        if (response.ok) {
            alert("Campaign deleted!");
            fetchMyCampaigns();
        } else {
            console.error("Error deleting campaign", response.status);
        }
    }

