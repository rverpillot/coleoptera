
$(document).on("element-created element-updated", function (event) {
    // console.log("element-created element-updated")
    doSemanticUI($("#" + event.detail.element))
})


var map_individus;
var gmarkers = [];

function mapIGN(map, controls) {
    L.geoportalLayer.WMTS({
        layer: "ORTHOIMAGERY.ORTHOPHOTOS"
    }).addTo(map);
    // L.geoportalLayer.WMTS({
    //     layer: "GEOGRAPHICALGRIDSYSTEMS.MAPS"
    // }).addTo(map);
    L.geoportalLayer.WMTS({
        layer: "GEOGRAPHICALGRIDSYSTEMS.PLANIGNV2"
    }, {
        opacity: 100
    }).addTo(map);

    var layerSwitcher = L.geoportalControl.LayerSwitcher();
    map.addControl(layerSwitcher);

    if (controls) {
        var mp = L.geoportalControl.MousePosition({
            displayCoordinate: true,
            altitude: {
                triggerDelay: 1000,
            }
        });
        map.addControl(mp);

        var search = L.geoportalControl.SearchEngine({ displayMarker: false });
        map.addControl(search);
    }
    return map
}

function showMarkers(tag, markers) {
    // console.log("show markers")

    $.each(gmarkers, function (i, marker) {
        marker.remove()
    })

    gmarkers = []
    positions = []
    $.each(markers, function (i, location) {
        var position = location.Location
        var marker = L.marker(position, {
            title: location.Infos.join("\n")
        })
        positions.push(position)
        gmarkers.push(marker)
        marker.addTo(map_individus)
    })

    var bounds = L.latLngBounds(positions)
    // map_individus.fitBounds(bounds, { maxZoom: 6, padding: [5, 5] })
    map_individus.fitBounds(bounds, { maxZoom: 6 })
}

function createMap(tag, center, zoom, pageName) {
    if ($(tag).length == 0) return;

    map_individus = L.map($(tag)[0], {})
    mapIGN(map_individus, true)

    if (center.lat == 0 && center.lng == 0) {
        map_individus.locate({ setView: true, maxZoom: 10 })
    } else {
        map_individus.setView([center.lat, center.lng], zoom)
    }

    map_individus.on("moveend zoomend", function (ev) {
        var center = map_individus.getCenter()
        var data = { lat: center.lat, lng: center.lng, zoom: map_individus.getZoom() }
        ihui.trigger("map-changed", pageName, data, false)
    })

    console.log("createMap")
}

function refreshMap(center, zoom) {
    console.log("refreshMap")
    //    map_individus.setView(center, zoom)
    map_individus.invalidateSize()
}


function createEditMap(tag, pageName) {
    var latitude = parseFloat($("[name=latitude]").val())
    var longitude = parseFloat($("[name=longitude]").val())
    var position = { lat: 47.626951, lng: 6.997541 }

    if (!isNaN(longitude) && !isNaN(latitude)) {
        position = { lat: latitude, lng: longitude };
    }
    var editMap = L.map($(tag)[0], {
        center: position,
        zoom: 10,
        scrollWheelZoom: false
    })
    mapIGN(editMap, true);

    var editMarker = L.marker(position, {
        // title: position.toString(),
        draggable: true
    });
    editMarker.addTo(editMap);
    // console.log("createEditMap")

    editMarker.on('dragend', function (event) {
        var position = editMarker.getLatLng()
        // editMarker.setTitle(position.toString())
        ihui.trigger("position", pageName, position)
    })

    editMap.on("moveend", function (ev) {
        var pos = editMap.getCenter()
        editMarker.setLatLng(pos)
    })

}

function createPreviewMap(tag, longitude, latitude) {
    var position = { lat: latitude, lng: longitude };

    var previewMap = L.map($(tag)[0], {
        center: position,
        zoom: 10,
        scrollWheelZoom: false
    })
    mapIGN(previewMap, false);

    var marker = L.marker(position, {})
    marker.addTo(previewMap)
    // console.log("createPreviewMap")
}


function doSemanticUI(tag) {
    $(tag).find('.ui.modal').modal({ closable: false }).modal("show")
    $(tag).find('.ui.dropdown').dropdown({
        forceSelection: false,
        fullTextSearch: true,
        message: {
            addResult: "Ajouter <b>{term}</b>",
            count: "{count} selectionné(s)",
            maxSelections: 'Max {maxCount} selections',
            noResults: 'Aucun résultat.'
        }
    })
    $(tag).find('select.dropdown').dropdown()

}
