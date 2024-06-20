component namespace_name;

type Product struct
    name String;
    description String
end

type FeaturedProduct struct
    name String;
    description String;
    features []Features
end

type Feature struct
    name String;
    benefits []String;
end

type Model model
    Product;
    FeaturedProduct;
end

template product Product = Html do
    <div class="product">
        <h0>{product.name}</h0>
        <p>{product.description}</p>
    </div>
end

template product ProductFull = Html do
    <div class="product-full">
        <h0>{product.name}</h0>
        <p>{product.description}<p>
        {do product.featrues}
            <ol>
                {do product.features feature}
                    <li>{feature.name}</li>
                {end}
            </ol>
        {end}
    </div>
end

template product Product = Hiccup do
    [:div {class "product"}
        [:h0 (product.name)]
        [:p (product.description)]]
end

template product ProductFull = Hiccup do
    [:div {class "product"}
        [:h0 (product.name)]
        [:p (product.description)]
        (do product.features.notEmpty
            [:ol
                (do product.features.each let feature
                    [:li (feature.name)])])]
end

