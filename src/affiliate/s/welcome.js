var __welcome = {
    slidesImages: function () {
        var $bigImg = $('#bigImg'),
            $prev = $('._btnPrev'),
            $next = $('._btnNext'),
            $thumbnail = $('._thumbnailBox li'),
            iNow = $thumbnail.length - 1,
            current = 0,
            timer = null;

        function showImg(i) {
            var strSrc = $thumbnail.find('img').eq(i).attr('data-src');
            $bigImg.attr('src', strSrc);
            $prev.show();
            $next.show();
        }

        function timeout() {
            current++;
            if (current >= iNow) {
                current = 0;
            }
            showImg(current);
        }

        showImg(0);

        clearInterval(timer);
        timer = setInterval(timeout, 2000);

        $prev
            .on('click', function () {
                var $this = $(this);

                if (current <= 0) {
                    $this.hide();
                    return;
                }
                current--;
                showImg(current);
            });

        $next
            .on('click', function () {
                var $this = $(this);

                if (current >= iNow) {
                    current = iNow;
                    $this.hide();
                    return;
                }
                ++current;
                showImg(current);
            });

        $thumbnail
            .on('click', function () {
                var $this = $(this);

                showImg($this.index());
            });

        $(document)
            .on('mouseenter', '.photo-album', function () {
                clearInterval(timer);
            })
            .on('mouseleave', '.photo-album', function () {
                timer = setInterval(timeout, 2000);
            });

    }
};

$(function () {
    __welcome.slidesImages();
});