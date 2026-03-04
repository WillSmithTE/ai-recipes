(function () {
    "use strict";

    var toggle = document.getElementById("feedback-toggle");
    var panel = document.getElementById("feedback-panel");
    var closeBtn = document.getElementById("feedback-close");
    var form = document.getElementById("feedback-form");
    var submitBtn = form.querySelector(".feedback-submit");
    var stars = form.querySelectorAll(".star-btn");
    var textarea = document.getElementById("feedback-comment");
    var successEl = document.getElementById("feedback-success");

    var selectedRating = 0;

    // Toggle panel open/close
    toggle.addEventListener("click", function () {
        var isOpen = panel.classList.contains("open");
        panel.classList.toggle("open");
        panel.setAttribute("aria-hidden", isOpen ? "true" : "false");
    });

    closeBtn.addEventListener("click", function () {
        panel.classList.remove("open");
        panel.setAttribute("aria-hidden", "true");
    });

    // Star rating selection
    stars.forEach(function (star) {
        star.addEventListener("click", function () {
            selectedRating = parseInt(star.dataset.rating, 10);
            updateStars();
            submitBtn.disabled = false;
        });
    });

    function updateStars() {
        stars.forEach(function (star) {
            var rating = parseInt(star.dataset.rating, 10);
            if (rating <= selectedRating) {
                star.classList.add("active");
            } else {
                star.classList.remove("active");
            }
        });
    }

    // Submit feedback
    form.addEventListener("submit", function (e) {
        e.preventDefault();

        var payload = {
            rating: selectedRating,
            comment: textarea.value.trim(),
            page: window.location.pathname,
        };

        submitBtn.disabled = true;
        submitBtn.textContent = "Sending...";

        fetch("/api/feedback", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })
            .then(function (res) {
                if (!res.ok) throw new Error("Request failed");
                return res.json();
            })
            .then(function () {
                form.hidden = true;
                successEl.hidden = false;

                // Auto-close after a delay
                setTimeout(function () {
                    panel.classList.remove("open");
                    panel.setAttribute("aria-hidden", "true");

                    // Reset form state after panel closes
                    setTimeout(function () {
                        form.hidden = false;
                        successEl.hidden = true;
                        selectedRating = 0;
                        textarea.value = "";
                        submitBtn.textContent = "Submit";
                        submitBtn.disabled = true;
                        updateStars();
                    }, 300);
                }, 2000);
            })
            .catch(function () {
                submitBtn.disabled = false;
                submitBtn.textContent = "Submit";
                alert("Failed to send feedback. Please try again.");
            });
    });
})();
