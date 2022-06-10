

// Hide the checkboxes.
// We only need to see them when "filters" has been pressed.
$(document).ready(function () {
    $('div.checkboxes').hide();
});

// Toggle hide/visible when pressing the filter button

$(document).ready(function () {
    $("p").click(function () {
        if ($('div.checkboxes').is(':visible')) {
            $("div.checkboxes").slideToggle()
            var boxes = []
            $("input:checkbox:checked").each(function() {
                boxes.push($(this).attr("value"))
            })
            $.post('/feed/', {"lang": boxes}, function(data) {

                var articles = $("div.article", "div.article_list")
                for (a in data) {
                    new_headline = '<a href="' + data[a].url + '">▤ ' + data[a].headline + "/>"
                    new_story = data[a].source + ": " + data[a].headline
                    new_publish_date = "Published: " + data[a].docdate
                    console.log(new_publish_date)

                    articles.eq(a).attr("id", data[a].id)
                    articles.eq(a).find("div.headline").html(new_headline)
                    articles.eq(a).find("div.story").html(new_story)
                    articles.eq(a).find("div.date").html(new_publish_date)
                }
            })
            //$("#filter").submit();
        } else {
            $("div.checkboxes").slideToggle();
        }
    });
});
